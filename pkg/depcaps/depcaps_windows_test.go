package depcaps_test

import "github.com/breml/depcaps/pkg/depcaps"

func osSpecificLinterSettings(linterSettings *depcaps.LinterSettings) *depcaps.LinterSettings {
	if linterSettings == nil {
		linterSettings = &depcaps.LinterSettings{}
	}
	if linterSettings.GlobalAllowedCapabilities == nil {
		linterSettings.GlobalAllowedCapabilities = make(map[string]bool, 1)
	}
	linterSettings.GlobalAllowedCapabilities["CAPABILITY_SYSTEM_CALLS"] = true

	return linterSettings
}
