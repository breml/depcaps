package depcaps_test

import (
	"testing"

	"github.com/breml/depcaps/pkg/depcaps"
)

func TestLinterSettingsSet(t *testing.T) {
	settings := &depcaps.LinterSettings{}
	err := settings.Set("testdata/ok.json")
	if err != nil {
		t.Fatalf("Failed to set settings: %s", err)
	}
	if settings.GlobalAllowedCapabilities["CAPABILITY_UNSPECIFIED"] != true {
		t.Fatalf("CAPABILITY_UNSPECIFIED not set")
	}
	if settings.PackageAllowedCapabilities["github.com/google/uuid"] == nil {
		t.Fatalf("github.com/google/uuid")
	}
	if settings.PackageAllowedCapabilities["github.com/google/uuid"]["CAPABILITY_RUNTIME"] != true {
		t.Fatalf("CAPABILITY_RUNTIME not set")
	}
}

func TestLinterSettingsSetError(t *testing.T) {
	tt := []struct {
		name     string
		filename string

		wantErr bool
	}{
		{
			name:     "file not found",
			filename: "testdata/notfound.json",
			wantErr:  true,
		},
		{
			name:     "invalid json",
			filename: "testdata/invalid.json",
			wantErr:  true,
		},
		{
			name:     "invalid global capability",
			filename: "testdata/invalid_global_capability.json",
			wantErr:  true,
		},
		{
			name:     "invalid package capability",
			filename: "testdata/invalid_package_capability.json",
			wantErr:  true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			settings := &depcaps.LinterSettings{}
			err := settings.Set(tc.filename)
			if err != nil && !tc.wantErr {
				t.Fatalf("Failed to set settings: %s", err)
			}
			if err == nil && tc.wantErr {
				t.Fatalf("Expected error, got none")
			}
		})
	}
}
