//go:build go1.22
// +build go1.22

package main

import (
	"github.com/google/uuid" // want "Package github.com/google/uuid has not allowed capability CAPABILITY_NETWORK" "Package github.com/google/uuid has not allowed capability CAPABILITY_REFLECT"
)

func main() {
	uuid.GetTime()
}
