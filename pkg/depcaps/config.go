package depcaps

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/capslock/proto"
)

type LinterSettings struct {
	GlobalAllowedCapabilities  map[string]bool            `json:"GlobalAllowedCapabilities"`
	PackageAllowedCapabilities map[string]map[string]bool `json:"PackageAllowedCapabilities"`
	CapslockBaselineFile       string                     `json:"-"`
}

func (s LinterSettings) IsBoolFlag() bool { return false }
func (s LinterSettings) Get() interface{} { return s }

func (s LinterSettings) String() string {
	return ""
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

	for c := range s.GlobalAllowedCapabilities {
		if _, ok := proto.Capability_value[c]; !ok {
			return fmt.Errorf("invalid global capability: %s", c)
		}
	}

	for p, pv := range s.PackageAllowedCapabilities {
		for c := range pv {
			if _, ok := proto.Capability_value[c]; !ok {
				return fmt.Errorf("invalid capability for package %q: %s", p, c)
			}
		}
	}

	return nil
}
