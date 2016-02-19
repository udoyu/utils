package timer

import (
	"sync"
	"time"
)

type TimerHandler struct {
	sync.RWMutex
	timermap map[time.Duration]*TimeWheel
}

func NewTimerHandler() *TimerHandler {
	return &TimerHandler{
		timermap: make(map[time.Duration]*TimeWheel),
	}
}

func (this *TimerHandler) get(d time.Duration) *TimeWheel {
	var tw *TimeWheel
	this.RLock()
	tw, _ = this.timermap[d]
	this.RUnlock()
	return tw
}

func (this *TimerHandler) add(d time.Duration) *TimeWheel {
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
	tw := NewTimeWheel(p, int(i))
	this.Lock()
	this.timermap[d] = tw
	this.Unlock()
	return tw
}

func (this *TimerHandler) AfterFunc(d time.Duration, task func()) {
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

func (this *TimerHandler) After(d time.Duration) <-chan struct{} {
	if d == 0 {
		return nil
	}
	tw := this.get(d)
	if tw == nil {
		tw = this.add(d)
	}
	return tw.After()
}

func (this *TimerHandler) Stop() {
	this.RLock()
	for _, v := range this.timermap {
		v.Stop()
	}
	defer this.RUnlock()
}
