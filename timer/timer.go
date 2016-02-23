package timer

import (
	"time"
)

var (
	globalPrecision = time.Millisecond * 100
	globalBaseMask  = 10
	globalxtimer    = NewXTimerHandler(globalBaseMask)
)

func ResetPrecision(precision time.Duration, stop ...bool) {
	globalPrecision = precision
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(globalBaseMask)
	if len(stop) > 0 && stop[0] {
		oldxtimer.Stop()
	}
}

func ResetMask(mask int, stop ...bool) {
	globalBaseMask = mask
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(globalBaseMask)
	if len(stop) > 0 && stop[0] {
		oldxtimer.Stop()
	}
}

func ResetPrecisionAndMask(precision time.Duration, mask int,stop ...bool) {
	globalPrecision = precision
	globalBaseMask = mask
	oldxtimer := globalxtimer
	globalxtimer = NewXTimerHandler(globalBaseMask)
	if len(stop) > 0 && stop[0] {
		oldxtimer.Stop()
	}
}

func After(d time.Duration) <-chan struct{} {
	return globalxtimer.After(d)
}

func AfterFunc(d time.Duration, f func()) {
	globalxtimer.AfterFunc(d, f)
}
