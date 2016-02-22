package timer

import (
	"time"
)

//1s = 100ms *10
//1min = 100ms * 100 * 6
//1hour = 3.60 * 100 * 100 * 100ms

var (
	ELEMENT_CNT_PER_BUCKET = []int64{256, 64, 64, 64, 64}
)

type Timer struct {
	timerWheels *XTimeWheel
	baseTime    time.Duration
}

func (this *Timer) init() {
	this.timerWheels = NewXTimeWheel(this.baseTime, ELEMENT_CNT_PER_BUCKET)
}

func NewTimer(baseTime ...time.Duration) *Timer {
	var d time.Duration
	if len(baseTime) > 0 {
		d = baseTime[0]
	} else {
		d = globalBaseTime
	}
	th := &Timer{
		baseTime: d,
	}
	th.init()
	return th
}

func (this *Timer) AfterFunc(d time.Duration, task func()) {
	this.timerWheels.AfterFunc(d, task)
}

func (this *Timer) After(d time.Duration) <-chan struct{} {
	return this.timerWheels.After(d)
}

func (this *Timer) Stop() {
	this.timerWheels.Stop()
}
