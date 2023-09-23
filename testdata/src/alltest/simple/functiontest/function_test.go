package functiontest

import (
	"testing"

	"github.com/google/uuid"
)

func TestCall(t *testing.T) {
	uuid.GetTime() // regular function call in test is not reported
}
