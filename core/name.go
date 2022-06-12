package core

import "go.uber.org/atomic"

var globalName = *atomic.NewString("INVALID")

func SetName(value string) {
	globalName.Store(value)
}

func Name() string {
	return globalName.Load()
}
