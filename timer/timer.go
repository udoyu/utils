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
	if d >= time.Second && d < time.Hour {
		p = time.Second
		i = int64(d) / int64(p)
	} else if d < time.Second && d >= time.Millisecond {
		p = time.Millisecond
		i = int64(d) / int64(p)
	} else if d < time.Millisecond && d >= time.Microsecond {
		p = time.Microsecond
		i = int64(d) / int64(p)
	} else if d < time.Microsecond {
		p = time.Nanosecond
		i = int64(d) / int64(p)
	} else if d >= time.Hour {
		p = time.Hour
		i = int64(d) / int64(p)
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

func (this *TimerHandler) Add(d time.Duration, task func()) {
	tw := this.get(d)
	if tw == nil {
		tw = this.add(d)
	}
	if tw == nil {
		go task()
	} else {
		tw.Add(task)
	}
}
