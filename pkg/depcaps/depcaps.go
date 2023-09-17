package depcaps

import (
	"flag"
	"fmt"
	"go/token"
	"strings"

	"github.com/google/capslock/analyzer"
	"github.com/google/capslock/proto"
	"golang.org/x/tools/go/analysis"

	"github.com/breml/depcaps/pkg/module"
)

const (
	doc = "depcaps maps capabilities of dependencies agains a set of allowed capabilities"
)

type depcaps struct {
	globalAllowedCapabilities  map[proto.Capability]struct{}
	packageAllowedCapabilities map[string]map[proto.Capability]struct{}
}

// NewAnalyzer returns a new depcaps analyzer.
func NewAnalyzer() *analysis.Analyzer {
	depcaps := depcaps{
		globalAllowedCapabilities:  make(map[proto.Capability]struct{}),
		packageAllowedCapabilities: make(map[string]map[proto.Capability]struct{}),
	}

	// TODO: read from config file or CLI argument
	// depcaps.globalAllowedCapabilities[proto.Capability_CAPABILITY_ARBITRARY_EXECUTION] = struct{}{}
	// depcaps.packageAllowedCapabilities["github.com/google/capslock/proto"] = make(map[proto.Capability]struct{})
	// depcaps.packageAllowedCapabilities["github.com/google/capslock/proto"][proto.Capability_CAPABILITY_SYSTEM_CALLS] = struct{}{}

	a := &analysis.Analyzer{
		Name: "depcaps",
		Doc:  doc,
		Run:  depcaps.run,
	}

	a.Flags.Init("depcaps", flag.ExitOnError)
	a.Flags.Var(versionFlag{}, "V", "print version and exit")

	return a
}

func (d depcaps) run(pass *analysis.Pass) (interface{}, error) {
	if isTestPackage(pass) {
		return nil, nil
	}

	packagePrefix := pass.Pkg.Path()

	moduleFile, err := module.GetModuleFile()
	if err == nil {
		packagePrefix = moduleFile.Module.Mod.Path
	}

	packageNames := []string{pass.Pkg.Path()}

	// TODO: decide on unanalyzed capabilities
	// TODO: is it possible to create a middleware / wrapper for a classifier?
	classifier := analyzer.GetClassifier(true)

	pkgs := analyzer.LoadPackages(packageNames,
		analyzer.LoadConfig{
			// TODO: support BuildTags, GOOS and GOARCH?
			// 	BuildTags: *buildTags,
			// 	GOOS:      *goos,
			// 	GOARCH:    *goarch,
		},
	)
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("No packages matching %v", packageNames)
	}

	queriedPackages := analyzer.GetQueriedPackages(pkgs)

	offendingCapabilities := make(map[string]map[proto.Capability]struct{})

	cil := analyzer.GetCapabilityInfo(pkgs, queriedPackages, classifier)
	for _, c := range cil.CapabilityInfo {
		if c.GetCapabilityType() != proto.CapabilityType_CAPABILITY_TYPE_TRANSITIVE {
			continue
		}

		if len(c.GetPath()) < 2 {
			panic("for transitive capabilities, a min length of 2 is expected")
		}

		pathName := *c.GetPath()[1].Name
		if strings.HasPrefix(pathName, packagePrefix) {
			// if we call an other package of our own module, we ignore this call here
			// TODO: make this behavior configurable
			continue
		}
		pkg := (pathName)[:strings.LastIndex(pathName, ".")]

		if len(pkg) == 0 {
			continue
		}

		if _, ok := d.globalAllowedCapabilities[c.GetCapability()]; ok {
			continue
		}
		if pkgAllowedCaps, ok := d.packageAllowedCapabilities[pkg]; ok {
			if _, ok := pkgAllowedCaps[c.GetCapability()]; ok {
				continue
			}
		}

		if _, ok := offendingCapabilities[pkg]; !ok {
			offendingCapabilities[pkg] = make(map[proto.Capability]struct{})
		}

		offendingCapabilities[pkg][c.GetCapability()] = struct{}{}
	}

	for pkg, pkgCaps := range offendingCapabilities {
		for cap := range pkgCaps {
			pos := findPos(pass, pkg)
			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				Message: fmt.Sprintf("Package %s has not allowed capability %s", pkg, cap),
			})
		}
	}

	return nil, nil
}

func isTestPackage(pass *analysis.Pass) bool {
	if strings.HasSuffix(pass.Pkg.Path(), ".test") || strings.HasSuffix(pass.Pkg.Path(), "_test") {
		return true
	}

	for _, f := range pass.Files {
		if strings.HasSuffix(pass.Fset.File(f.Pos()).Name(), "_test.go") {
			return true
		}
	}

	return false
}

func findPos(pass *analysis.Pass, pkg string) token.Pos {
	for _, file := range pass.Files {
		for _, i := range file.Imports {
			if pkg == strings.Trim(i.Path.Value, `"`) {
				return i.Pos()
			}
		}
	}

	return token.NoPos
}
