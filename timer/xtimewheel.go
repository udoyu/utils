package timer

import (
	"sync/atomic"
	"time"
)

type XTimeWheel struct {
	timers []*TimeWheel
	index  uint32
	size   uint32
}

func NewXTimeWheel(precision time.Duration, interval, size int) *XTimeWheel {
	xt := new(XTimeWheel)
	xt.timers = make([]*TimeWheel, size)
	xt.size = uint32(size)
	for i := 0; i < size; i++ {
		xt.timers[i] = NewTimeWheel(precision, interval)
	}
	return xt
}

func (this *XTimeWheel) AfterFunc(task func()) {
	i := atomic.AddUint32(&this.index, 1) % this.size
	this.timers[i].AfterFunc(task)
}

func (this *XTimeWheel) After() <-chan struct{} {
	return this.timers[0].After()
}

func (this *XTimeWheel) Stop() {
	for _, t := range this.timers {
		t.Stop()
	}
}
