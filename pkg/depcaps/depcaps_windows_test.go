package depcaps_test

import "github.com/breml/depcaps/pkg/depcaps"

func osSpecificLinterSettings(linterSettings *depcaps.LinterSettings) {
	linterSettings.GlobalAllowedCapabilities["CAPABILITY_SYSTEM_CALLS"] = true
}
