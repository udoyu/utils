package timer

import (
	"time"
)

//1s = 100ms *10
//1min = 100ms * 100 * 6
//1hour = 3.60 * 100 * 100 * 100ms

var (
	globalBaseTime         = time.Millisecond * 100
	BUCKET_CNT             = 5
	ELEMENT_CNT_PER_BUCKET = []int64{256, 64, 64, 64, 64}
	RIGHT_SHIFT_PER_BUCKET = []int64{8, 6, 6, 6, 6}
	BASE_PER_BUCKET        = []int64{1, 256, 256 * 64, 256 * 64 * 64, 256 * 64 * 64 * 64}
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

//func (this *TimerHandler) get(d time.Duration) (int, bool) {
//	nd := int64(d / this.baseTime)

//	if nd < BASE_PER_BUCKET[0] {
//		return 0, false
//	} else if nd < BASE_PER_BUCKET[1] {
//		return 0, true
//	} else if nd < BASE_PER_BUCKET[2] {
//		return 1, true
//	} else if nd < BASE_PER_BUCKET[3] {
//		return 2, true
//	} else if nd < BASE_PER_BUCKET[4] {
//		return 3, true
//	} else {
//		return 4, true
//	}
//	return 0, false
//}

func (this *Timer) AfterFunc(d time.Duration, task func()) {
	this.timerWheels.AfterFunc(d, task)
}

func (this *Timer) After(d time.Duration) <-chan struct{} {
	return this.timerWheels.After(d)
}

func (this *Timer) Stop() {
	this.timerWheels.Stop()
}
