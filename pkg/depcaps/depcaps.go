package depcaps

import (
	"flag"
	"fmt"
	"go/token"
	"os"
	"strings"
	"sync"

	"github.com/google/capslock/analyzer"
	"github.com/google/capslock/proto"
	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/breml/depcaps/pkg/module"
)

type depcaps struct {
	*LinterSettings
}

// NewAnalyzer returns a new depcaps analyzer.
func NewAnalyzer(settings *LinterSettings) *analysis.Analyzer {
	depcaps := depcaps{
		LinterSettings: &LinterSettings{
			GlobalAllowedCapabilities:  map[string]bool{},
			PackageAllowedCapabilities: map[string]map[string]bool{},
		},
	}

	a := &analysis.Analyzer{
		Name:     "depcaps",
		Doc:      "depcaps maps capabilities of dependencies agains a set of allowed capabilities",
		Run:      depcaps.run,
		Requires: []*analysis.Analyzer{},
	}

	a.Flags.Init("depcaps", flag.ExitOnError)
	a.Flags.Var(versionFlag{}, "V", "print version and exit")
	a.Flags.Var(depcaps.LinterSettings, "config", "depcaps linter settings config file")
	a.Flags.StringVar(&depcaps.CapslockBaselineFile, "reference", "", "capslock capabilities reference file")

	// if settings are provided, these have precedence
	if settings != nil {
		depcaps.GlobalAllowedCapabilities = settings.GlobalAllowedCapabilities
		depcaps.PackageAllowedCapabilities = settings.PackageAllowedCapabilities
		depcaps.CapslockBaselineFile = settings.CapslockBaselineFile
	}

	return a
}

var (
	once       = sync.Once{}
	mu         sync.Mutex
	stdSet     = make(map[string]struct{})
	moduleFile *modfile.File
	cil        *proto.CapabilityInfoList
	baseline   *proto.CapabilityInfoList

	osArgs []string
)

func (d *depcaps) run(pass *analysis.Pass) (interface{}, error) {
	// init depcaps linter
	{
		// TODO: can the `WithContextSetter` from golangci-lint help here?
		var err error
		once.Do(func() {
			mu.Lock()
			defer mu.Unlock()

			// init std pkg list
			var stdPkgs []*packages.Package
			stdPkgs, err = packages.Load(&packages.Config{Tests: false}, "std")
			if err != nil {
				return // error is returned after the once.Do-block
			}

			pre := func(pkg *packages.Package) bool {
				stdSet[pkg.PkgPath] = struct{}{}
				return true
			}
			packages.Visit(stdPkgs, pre, nil)

			// init moduleFile
			moduleFile, err = module.GetModuleFile()
			if err != nil {
				return // err is returned after the once.Do-block
			}

			packageNames := flag.Args()
			if len(packageNames) == 0 {
				packageNames = []string{"."}

				if len(osArgs) > 0 {
					packageNames = osArgs
				}
			}

			var classifier analyzer.Classifier = analyzer.GetClassifier(true)

			pkgs := analyzer.LoadPackages(packageNames,
				analyzer.LoadConfig{
					// TODO: support BuildTags, GOOS and GOARCH?
					// 	BuildTags: *buildTags,
					// 	GOOS:      *goos,
					// 	GOARCH:    *goarch,
				},
			)
			if len(pkgs) == 0 {
				err = fmt.Errorf("no packages matching %v", packageNames)
				return // err is returned after the once.Do.block
			}

			queriedPackages := analyzer.GetQueriedPackages(pkgs)
			cil = analyzer.GetCapabilityInfo(pkgs, queriedPackages, classifier)

			err := readCapslockBaseline(d.CapslockBaselineFile)
			if err != nil {
				return // err is returned after the once.Do.block
			}
		})
		if err != nil { // process error from packages.Load if executed once and it returned an error
			return nil, err
		}
	}

	if isTestPackage(pass) {
		return nil, nil
	}

	packageName := pass.Pkg.Path()
	packagePrefix := pass.Pkg.Path()
	if moduleFile != nil {
		packagePrefix = getModulePath()
	}

	mu.Lock()
	defer mu.Unlock()

	offendingCapabilities := make(map[string]map[proto.Capability]struct{})
	if baseline != nil {
		offendingCapabilities = diffCapabilityInfoLists(baseline, cil, packageName, packagePrefix)
	}

	for _, ci := range cil.GetCapabilityInfo() {
		depPkg, skip := relevantCapabilityInfo(ci, packageName, packagePrefix)
		if !skip {
			continue
		}

		if _, ok := stdSet[depPkg]; ok {
			continue
		}

		if _, ok := offendingCapabilities[depPkg]; !ok {
			offendingCapabilities[depPkg] = make(map[proto.Capability]struct{})
		}

		if ok := d.GlobalAllowedCapabilities[ci.Capability.String()]; ok {
			delete(offendingCapabilities[depPkg], ci.GetCapability())
			continue
		}
		if pkgAllowedCaps, ok := d.PackageAllowedCapabilities[depPkg]; ok {
			if ok := pkgAllowedCaps[ci.Capability.String()]; ok {
				delete(offendingCapabilities[depPkg], ci.GetCapability())
				continue
			}
		}
		if baseline != nil {
			continue
		}

		offendingCapabilities[depPkg][ci.GetCapability()] = struct{}{}
	}

	// TODO: sort offendingCapabilities by package name and capability name before reporting
	for pkg, pkgCaps := range offendingCapabilities {
		for cap := range pkgCaps {
			pos := findPos(pass, pkg)
			if pos == 0 {
				// TODO: figure out, if and why this is necessary
				// skip offending capabilites, that can not be matched to a source code location
				continue
			}

			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: fmt.Sprintf("Package %s has not allowed capability %s", pkg, cap),
			})
		}
	}

	return nil, nil
}

func readCapslockBaseline(capslockBaselineFile string) error {
	if capslockBaselineFile == "" {
		return nil
	}

	baselineData, err := os.ReadFile(capslockBaselineFile)
	if err != nil {
		return fmt.Errorf("Error reading baseline file: %v", err)
	}
	baseline = &proto.CapabilityInfoList{}
	err = protojson.Unmarshal(baselineData, baseline)
	if err != nil {
		return fmt.Errorf("Baseline file should include output from running `capslock -output=j`. Error parsing baseline file: %v", err)
	}
	return nil
}

func getModulePath() string {
	mu.Lock()
	defer mu.Unlock()

	return moduleFile.Module.Mod.Path
}

func isTestPackage(pass *analysis.Pass) bool {
	if strings.HasSuffix(pass.Pkg.Path(), ".test") || strings.HasSuffix(pass.Pkg.Path(), "_test") {
		return true
	}

	// for _, f := range pass.Files {
	// 	if strings.HasSuffix(pass.Fset.File(f.Pos()).Name(), "_test.go") {
	// 		return true
	// 	}
	// }

	return false
}

func relevantCapabilityInfo(ci *proto.CapabilityInfo, packageName, packagePrefix string) (string, bool) {
	if ci.GetCapabilityType() != proto.CapabilityType_CAPABILITY_TYPE_TRANSITIVE {
		return "", false
	}

	if len(ci.GetPath()) < 2 {
		panic("for transitive capabilities, a min length of 2 is expected")
	}

	if extractPackagePath(*ci.GetPath()[0].Name) != packageName {
		return "", false
	}

	depPkg := extractPackagePath(*ci.GetPath()[1].Name)
	if strings.HasPrefix(depPkg, packagePrefix) {
		// if we call an other package of our own module, we ignore this call here
		// TODO: make this behavior configurable
		return "", false
	}

	if len(depPkg) == 0 {
		return "", false
	}

	return depPkg, true
}

func extractPackagePath(pathName string) string {
	if strings.HasPrefix(pathName, "(") { // method
		pathName = pathName[1:strings.LastIndex(pathName, ")")]
	}
	pathName = strings.TrimLeft(pathName, "*") // pointer receiver
	return pathName[:strings.LastIndex(pathName, ".")]
}

func findPos(pass *analysis.Pass, pkg string) token.Pos {
	for _, file := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(file.Pos()).Name(), "_test.go") {
			continue
		}
		for _, i := range file.Imports {
			if pkg == strings.Trim(i.Path.Value, `"`) {
				return i.Pos()
			}
		}
	}

	return token.NoPos
}
