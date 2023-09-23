package main

import (
	"alltest/simple/direct"
	"alltest/simple/function"
	"alltest/simple/method"
)

func main() {
	// calls to packages of the same module are not reported
	function.Call()
	method.Call()
	direct.Call()
}
