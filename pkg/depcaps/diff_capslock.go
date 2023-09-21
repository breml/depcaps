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

	"github.com/google/capslock/proto"
)

type (
	capabilitySet   map[proto.Capability]*proto.CapabilityInfo
	capabilitiesMap map[string]capabilitySet
)

// populateMap takes a CapabilityInfoList and returns a map from package
// directory and capability to a pointer to the corresponding entry in the
// input.
func populateMap(cil *proto.CapabilityInfoList, packageName, packagePrefix string) capabilitiesMap {
	m := make(capabilitiesMap)
	for _, ci := range cil.GetCapabilityInfo() {
		depPkg, skip := relevantCapabilityInfo(ci, packageName, packagePrefix)
		if !skip {
			continue
		}

		capmap := m[depPkg]
		if capmap == nil {
			capmap = make(capabilitySet)
			m[depPkg] = capmap
		}
		capmap[ci.GetCapability()] = ci
	}
	return m
}

func diffCapabilityInfoLists(baseline, current *proto.CapabilityInfoList, packageName, packagePrefix string) map[string]map[proto.Capability]struct{} {
	baselineMap := populateMap(baseline, packageName, packagePrefix)
	currentMap := populateMap(current, packageName, packagePrefix)

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
