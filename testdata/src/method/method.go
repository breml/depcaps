package main

import (
	"github.com/google/uuid" // want "Package github.com/google/uuid has not allowed capability CAPABILITY_REFLECT"
)

func main() {
	(&uuid.NullUUID{}).UnmarshalJSON(nil)
}
