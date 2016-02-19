package timer

import (
	"time"
)

var (
	globalxtimer = NewXTimerHandler(256)
)

func After(d time.Duration) <-chan struct{} {
	return globalxtimer.After(d)
}

func AfterFunc(d time.Duration, f func()) {
	globalxtimer.AfterFunc(d, f)
}