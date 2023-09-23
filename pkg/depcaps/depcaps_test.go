package depcaps_test

import (
	"os"
	"path/filepath"
	"testing"

	// unused import of "github.com/google/uuid" to workaround GOPROXY=no in
	// analysistest. This caches the module in CI before analysistest is executed.
	_ "github.com/google/uuid"
	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/breml/depcaps/pkg/depcaps"
)

func TestAll(t *testing.T) {
	tt := []struct {
		name           string
		linterSettings *depcaps.LinterSettings
		testdataDir    string
		packages       []string
	}{
		{
			name:           "init",
			linterSettings: nil,
			testdataDir:    "alltest",
			packages:       []string{"."},
		},
		{
			name:           "simple",
			linterSettings: nil,
			testdataDir:    "alltest",
			packages:       []string{"./simple/..."},
		},
		{
			name: "capslock file empty",
			linterSettings: &depcaps.LinterSettings{
				CapslockBaselineFile: "simple/capslock.json",
			},
			testdataDir: "alltest",
			packages:    []string{"./simple/..."},
		},
		{
			name: "global allow",
			linterSettings: &depcaps.LinterSettings{
				GlobalAllowedCapabilities: map[string]bool{
					"CAPABILITY_FILES": true,
				},
			},
			testdataDir: "alltest",
			packages:    []string{"./allow/..."},
		},
		{
			name: "package allow",
			linterSettings: &depcaps.LinterSettings{
				PackageAllowedCapabilities: map[string]map[string]bool{
					"github.com/google/uuid": {
						"CAPABILITY_FILES": true,
					},
				},
			},
			testdataDir: "alltest",
			packages:    []string{"./allow/..."},
		},
		{
			name: "capslock file empty with global allow",
			linterSettings: &depcaps.LinterSettings{
				GlobalAllowedCapabilities: map[string]bool{
					"CAPABILITY_FILES": true,
				},
				CapslockBaselineFile: "allow/capslock.json",
			},
			testdataDir: "alltest",
			packages:    []string{"./allow/..."},
		},
		{
			name: "capslock file empty with package allow",
			linterSettings: &depcaps.LinterSettings{
				PackageAllowedCapabilities: map[string]map[string]bool{
					"github.com/google/uuid": {
						"CAPABILITY_FILES": true,
					},
				},
				CapslockBaselineFile: "allow/capslock.json",
			},
			testdataDir: "alltest",
			packages:    []string{"./allow/..."},
		},
		{
			name: "capslock file",
			linterSettings: &depcaps.LinterSettings{
				CapslockBaselineFile: "capslockfile/capslock.json",
			},
			testdataDir: "alltest",
			packages:    []string{"./capslockfile/..."},
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

			tc.linterSettings = osSpecificLinterSettings(tc.linterSettings)
			depcaps.SetOSArgs([]string{"./..."})
			if tc.linterSettings != nil {
				depcaps.SetBaseline(tc.linterSettings.CapslockBaselineFile)
			}
			analysistest.Run(t, testCaseDir, depcaps.NewAnalyzer(tc.linterSettings), tc.packages...)
		})
	}
}
