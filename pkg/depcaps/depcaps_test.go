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
		name           string
		analyzerFunc   func(depcaps.LinterSettings) *analysis.Analyzer
		linterSettings depcaps.LinterSettings
		testdataDir    string
	}{
		{
			name:           "simple",
			analyzerFunc:   depcaps.NewAnalyzer,
			linterSettings: depcaps.LinterSettings{},
			testdataDir:    "simple",
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testCaseDir := filepath.Join(testdata, "src", tc.testdataDir)
			err = os.Chdir(testCaseDir)
			if err != nil {
				t.Fatalf("Failed to change wd: %s", err)
			}
			defer func() {
				err := os.Chdir(wd)
				if err != nil {
					t.Fatalf("Failed to return to wd: %s", err)
				}
			}()

			analysistest.Run(t, testCaseDir, tc.analyzerFunc(tc.linterSettings), ".")
		})
	}
}
