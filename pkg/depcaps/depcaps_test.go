package depcaps_test

import (
	"os"
	"path/filepath"
	"testing"

	// unused import of "github.com/google/uuid" to workaround GOPROXY=no in
	// analysistest. This caches the module in CI before analysistest is executed.
	_ "github.com/google/uuid"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/breml/depcaps/pkg/depcaps"
)

func TestAll(t *testing.T) {
	tt := []struct {
		name           string
		analyzerFunc   func(*depcaps.LinterSettings) *analysis.Analyzer
		linterSettings *depcaps.LinterSettings
		testdataDir    string
	}{
		{
			name:           "func",
			analyzerFunc:   depcaps.NewAnalyzer,
			linterSettings: nil,
			testdataDir:    "func",
		},
		{
			name:           "method",
			analyzerFunc:   depcaps.NewAnalyzer,
			linterSettings: nil,
			testdataDir:    "method",
		},
		{
			name:         "global allow",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				GlobalAllowedCapabilities: map[string]bool{
					"CAPABILITY_FILES": true,
				},
			},
			testdataDir: "allow",
		},
		{
			name:         "package allow",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				PackageAllowedCapabilities: map[string]map[string]bool{
					"github.com/google/uuid": {
						"CAPABILITY_FILES": true,
					},
				},
			},
			testdataDir: "allow",
		},
		{
			name:         "capslock file",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				CapslockBaselineFile: "capslock.json",
			},
			testdataDir: "capslockfile",
		},
		{
			name:         "capslock file empty",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				CapslockBaselineFile: "capslock.json",
			},
			testdataDir: "capslockfileempty",
		},
		{
			name:         "capslock file empty with global allow",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				GlobalAllowedCapabilities: map[string]bool{
					"CAPABILITY_FILES": true,
				},
				CapslockBaselineFile: "capslock.json",
			},
			testdataDir: "capslockfileemptyallow",
		},
		{
			name:         "capslock file empty with package allow",
			analyzerFunc: depcaps.NewAnalyzer,
			linterSettings: &depcaps.LinterSettings{
				PackageAllowedCapabilities: map[string]map[string]bool{
					"github.com/google/uuid": {
						"CAPABILITY_FILES": true,
					},
				},
				CapslockBaselineFile: "capslock.json",
			},
			testdataDir: "capslockfileemptyallow",
		},
		{
			name:           "stdlib",
			analyzerFunc:   depcaps.NewAnalyzer,
			linterSettings: nil,
			testdataDir:    "stdlib",
		},
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			depcaps.ResetGlobalState()

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

			tc.linterSettings = osSpecificLinterSettings(tc.linterSettings)

			analysistest.Run(t, testCaseDir, tc.analyzerFunc(tc.linterSettings), ".")
		})
	}
}
