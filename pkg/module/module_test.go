package module_test

import (
	"os"
	"testing"

	"github.com/breml/depcaps/pkg/module"
)

func TestGetModuleFile(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer func() {
		_ = os.Chdir(wd)
	}()

	err = os.Chdir("./testdata/a/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	file, err := module.GetModuleFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "github.com/breml/depcaps/testdata/a"
	if expected != file.Module.Mod.Path {
		t.Fatalf("expected %q, got: %q", expected, file.Module.Mod.Path)
	}
}

func TestGetModuleFile_here(t *testing.T) {
	file, err := module.GetModuleFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "github.com/breml/depcaps"
	if expected != file.Module.Mod.Path {
		t.Fatalf("expected %q, got %q", expected, file.Module.Mod.Path)
	}
}
