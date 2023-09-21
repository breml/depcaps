package depcaps

import "sync"

func ResetGlobalState() {
	mu.Lock()
	defer mu.Unlock()

	once = sync.Once{}
	initialized = false
	stdSet = make(map[string]struct{})
	moduleFile = nil
	cil = nil
}
