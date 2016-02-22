package timer

import (
	"sync/atomic"
	"time"
)

type XTimerHandler struct {
	handers []*Timer
	buckets int32
	index   int32
}

func NewXTimerHandler(buckets int) *XTimerHandler {
	xtw := &XTimerHandler{
		handers: make([]*Timer, buckets),
		buckets: int32(buckets),
	}
	for i := 0; i < buckets; i++ {
		xtw.handers[i] = NewTimer()
	}
	return xtw
}

func (this *XTimerHandler) After(d time.Duration) <-chan struct{} {
	index := atomic.AddInt32(&this.index, 1) % this.buckets
	return this.handers[index].After(d)
}

func (this *XTimerHandler) AfterFunc(d time.Duration, task func()) {
	index := atomic.AddInt32(&this.index, 1) % this.buckets
	this.handers[index].AfterFunc(d, task)
}

func (this *XTimerHandler) Stop() {
	for _, v := range this.handers {
		v.Stop()
	}
}
