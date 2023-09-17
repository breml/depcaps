package depcaps_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/breml/depcaps/pkg/depcaps"
)

func TestAll(t *testing.T) {
	tt := []struct {
		name         string
		analyzerFunc func() *analysis.Analyzer
		testdataDir  string
	}{
		{
			name:         "simple",
			analyzerFunc: depcaps.NewAnalyzer,
			testdataDir:  "simple",
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			analysistest.Run(t, testdata, tc.analyzerFunc(), tc.testdataDir)
		})
	}
}
