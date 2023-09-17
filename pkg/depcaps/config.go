package depcaps

import (
	"encoding/json"
	"os"
)

type LinterSettings struct {
	// TODO: Validate Capabilities
	GlobalAllowedCapabilities  map[string]bool            `json:"GlobalAllowedCapabilities"`
	PackageAllowedCapabilities map[string]map[string]bool `json:"PackageAllowedCapabilities"`
}

func (s LinterSettings) IsBoolFlag() bool { return false }
func (s LinterSettings) Get() interface{} { return s }

func (s LinterSettings) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *LinterSettings) Set(in string) error {
	b, err := os.ReadFile(in)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, s)
	if err != nil {
		return err
	}

	// TODO: Validate Capabilities

	return nil
}
