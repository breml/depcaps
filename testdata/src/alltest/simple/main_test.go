package main_test

import (
	"testing"

	"github.com/google/uuid"
)

func TestFunc(t *testing.T) {
	uuid.GetTime() // regular function call in tests is not reported
}
