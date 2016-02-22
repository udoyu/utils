package timer

import (
	"time"
)

var (
	globalBaseTime = time.Millisecond * 100
	globalBaseMask = 10
	globalxtimer   = NewXTimerHandler(globalBaseMask)
)

func ResetBaseTime(d time.Duration) {
	globalBaseTime = d
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(globalBaseMask)
	oldxtimer.Stop()
}

func ResetBaseMask(mask int) {
	globalBaseMask = mask
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(globalBaseMask)
	oldxtimer.Stop()
}

func After(d time.Duration) <-chan struct{} {
	return globalxtimer.After(d)
}

func AfterFunc(d time.Duration, f func()) {
	globalxtimer.AfterFunc(d, f)
}
