package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/breml/depcaps/pkg/depcaps"
)

var (
	version = "development"
	commit  = ""
	date    = ""
)

func main() {
	depcaps.Version = buildVersion()

	analyzer := depcaps.New(nil).WithFlagArgs().AsAnalyzer(true)

	singlechecker.Main(analyzer)
}

func buildVersion() string {
	result := fmt.Sprintf("%s version %s", filepath.Base(os.Args[0]), version)

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
		result = fmt.Sprintf("%s\nmodule version: %s, checksum: %s", result, info.Main.Version, info.Main.Sum)
	}
	result = fmt.Sprintf("%s\ngoos: %s\ngoarch: %s", result, runtime.GOOS, runtime.GOARCH)

	return result
}
