// Copyright 2023 Google LLC
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Originally based on code from:
// https://github.com/google/capslock/blob/dff452901025354c39a85dc4f942963920f0b519/analyzer/scan.go
// Modified to fit the needs of depcaps.

package depcaps

import (
	"sort"
	"strings"

	"github.com/google/capslock/proto"
)

type (
	capabilitySet   map[proto.Capability]*proto.CapabilityInfo
	capabilitiesMap map[string]capabilitySet
)

// populateMap takes a CapabilityInfoList and returns a map from package
// directory and capability to a pointer to the corresponding entry in the
// input.
func populateMap(cil *proto.CapabilityInfoList, packagePrefix string) capabilitiesMap {
	m := make(capabilitiesMap)
	// TODO: this loop is very similar to the loop in depcaps:run -> cleanup this code
	for _, ci := range cil.GetCapabilityInfo() {
		if ci.GetCapabilityType() != proto.CapabilityType_CAPABILITY_TYPE_TRANSITIVE {
			continue
		}

		if len(ci.GetPath()) < 2 {
			panic("for transitive capabilities, a min length of 2 is expected")
		}

		pathName := *ci.GetPath()[1].Name
		if strings.HasPrefix(pathName, packagePrefix) {
			// if we call an other package of our own module, we ignore this call here
			// TODO: make this behavior configurable
			continue
		}
		pkg := (pathName)[:strings.LastIndex(pathName, ".")]

		if len(pkg) == 0 {
			continue
		}

		capmap := m[pkg]
		if capmap == nil {
			capmap = make(capabilitySet)
			m[pkg] = capmap
		}
		capmap[ci.GetCapability()] = ci
	}
	return m
}

func diffCapabilityInfoLists(baseline, current *proto.CapabilityInfoList, packagePrefix string) map[string]map[proto.Capability]struct{} {
	baselineMap := populateMap(baseline, packagePrefix)
	currentMap := populateMap(current, packagePrefix)

	var packages []string
	for packageName := range baselineMap {
		packages = append(packages, packageName)
	}
	for packageName := range currentMap {
		if _, ok := baselineMap[packageName]; !ok {
			packages = append(packages, packageName)
		}
	}
	sort.Strings(packages)

	offendingCapabilities := make(map[string]map[proto.Capability]struct{})

	for _, packageName := range packages {
		if _, ok := offendingCapabilities[packageName]; !ok {
			offendingCapabilities[packageName] = make(map[proto.Capability]struct{})
		}
		b := baselineMap[packageName]
		c := currentMap[packageName]
		for capability := range c {
			if _, ok := b[capability]; !ok {
				offendingCapabilities[packageName][capability] = struct{}{}
			}
		}
	}

	return offendingCapabilities
}
