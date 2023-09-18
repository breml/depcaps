//go:build tools
// +build tools

package main

// for more about tools.go please visit
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/google/capslock/cmd/capslock"
)
