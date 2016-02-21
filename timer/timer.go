package timer

import (
	"time"
)

var (
	baseMask     = 10
	globalxtimer = NewXTimerHandler(baseMask)
)

func ResetBaseTime(d time.Duration) {
	globalBaseTime = d
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(baseMask)
	oldxtimer.Stop()
}

func ResetBaseMask(mask int) {
	baseMask = mask
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(baseMask)
	oldxtimer.Stop()
}

func After(d time.Duration) <-chan struct{} {
	return globalxtimer.After(d)
}

func AfterFunc(d time.Duration, f func()) {
	globalxtimer.AfterFunc(d, f)
}
