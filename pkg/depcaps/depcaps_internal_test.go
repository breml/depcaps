package depcaps

func SetOSArgs(args []string) {
	osArgs = args
}

func SetBaseline(baselineFile string) {
	mu.Lock()
	defer mu.Unlock()

	err := readCapslockBaseline(baselineFile)
	if err != nil {
		panic(err)
	}
}
