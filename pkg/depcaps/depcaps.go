package depcaps

import (
	"flag"

	"golang.org/x/tools/go/analysis"
)

const (
	doc = "depcaps maps capabilities of dependencies agains a set of allowed capabilities"
)

type depcaps struct{}

// NewAnalyzer returns a new depcaps analyzer.
func NewAnalyzer() *analysis.Analyzer {
	depcaps := depcaps{}

	a := &analysis.Analyzer{
		Name: "depcaps",
		Doc:  doc,
		Run:  depcaps.run,
	}

	a.Flags.Init("depcaps", flag.ExitOnError)
	a.Flags.Var(versionFlag{}, "V", "print version and exit")

	return a
}

func (b depcaps) run(pass *analysis.Pass) (interface{}, error) {
	return nil, nil
}
