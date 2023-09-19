//go:build !windows
// +build !windows

package depcaps_test

import (
	"github.com/breml/depcaps/pkg/depcaps"
)

func osSpecificLinterSettings(linterSettings *depcaps.LinterSettings) {}
