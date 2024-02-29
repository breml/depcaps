//go:build go1.22
// +build go1.22

package function

import (
	"github.com/google/uuid" // want "Package github.com/google/uuid has not allowed capability CAPABILITY_FILES" "Package github.com/google/uuid has not allowed capability CAPABILITY_NETWORK" "Package github.com/google/uuid has not allowed capability CAPABILITY_REFLECT"
)

func Call() {
	uuid.GetTime() // regular function call is reported
}
