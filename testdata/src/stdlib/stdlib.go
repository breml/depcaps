package main

import (
	"os"
)

func main() {
	os.Getenv("FOOBAR") // stdlib is not reported
}
