package function

import (
	"github.com/google/uuid" // want "Package github.com/google/uuid has not allowed capability CAPABILITY_FILES" "Package github.com/google/uuid has not allowed capability CAPABILITY_NETWORK"
)

func Call() {
	uuid.GetTime() // regular function call is reported
}
