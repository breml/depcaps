package direct

import (
	"os"
)

func Call() {
	os.Getenv("FOOBAR") // direct call to stdlib is not reported
}
