package timer

import (
	"sync"
	"time"
)

type XTimerHandler struct {
	sync.RWMutex
	timermap map[time.Duration]*XTimeWheel
	buckets  int
}

func NewXTimerHandler(buckets int) *XTimerHandler {
	return &XTimerHandler{
		timermap: make(map[time.Duration]*XTimeWheel),
		buckets : buckets,
	}
}

func (this *XTimerHandler) get(d time.Duration) *XTimeWheel {
	var tw *XTimeWheel
	this.RLock()
	tw, _ = this.timermap[d]
	this.RUnlock()
	return tw
}

func (this *XTimerHandler) add(d time.Duration) *XTimeWheel {
	p := time.Second
	i := int64(0)
	if d >= time.Microsecond {
		p = d / 10
		i = 10
	} else {
		p = time.Nanosecond
	}

	if i == 0 {
		return nil
	}
	tw := NewXTimeWheel(p, int(i), this.buckets)
	this.Lock()
	this.timermap[d] = tw
	this.Unlock()
	return tw
}

func (this *XTimerHandler) After(d time.Duration) <-chan struct{} {
	if d == 0 {
		return nil
	}
	tw := this.get(d)
	if tw == nil {
		tw = this.add(d)
	}
	return tw.After()
}

func (this *XTimerHandler) AfterFunc(d time.Duration, task func()) {
	if d == 0 {
		go task()
	}
	tw := this.get(d)
	if tw == nil {
		tw = this.add(d)
	}
	if tw == nil {
		go task()
	} else {
		tw.AfterFunc(task)
	}
}

func (this *XTimerHandler) Stop() {
	this.RLock()
	for _, v := range this.timermap {
		v.Stop()
	}
	defer this.RUnlock()
}
