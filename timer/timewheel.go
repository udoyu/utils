package timer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

type MutexList struct {
	sync.Mutex
	Elems list.List
}

type MutextTicker struct {
	c chan struct{}
	sync.RWMutex
}

func NewMutexTicker() MutextTicker {
	return MutextTicker{
		c: make(chan struct{}),
	}
}

func (this *MutextTicker) Get() <-chan struct{} {
	this.RLock()
	c := this.c
	this.RUnlock()
	return c
}

func (this *MutextTicker) Notify() {
	this.Lock()
	c := this.c
	this.c = make(chan struct{})
	this.Unlock()
	close(c)
}

type TimeWheel struct {
	ticket   *time.Ticker
	tasks    []MutexList
	interval int32
	curIndex int32
	status   int32
	timerC   []MutextTicker
}

func NewTimeWheel(precision time.Duration, interval int) *TimeWheel {
	tw := &TimeWheel{}
	tw.ticket = time.NewTicker(precision)
	tw.tasks = make([]MutexList, int(interval))
	for i := 0; i < int(interval); i++ {
		tw.tasks[i].Elems.Init()
	}
	tw.interval = int32(interval)
	tw.curIndex = int32(0)
	tw.timerC = make([]MutextTicker, interval)
	for i := 0; i < interval; i++ {
		tw.timerC[i] = NewMutexTicker()
	}
	go func(tw *TimeWheel) {
		defer tw.ticket.Stop()
		for atomic.LoadInt32(&tw.status) == 0 {
			select {
			case <-tw.ticket.C:
				tw.onTicker()
			}
		}
	}(tw)
	return tw
}

func (this *TimeWheel) After() <-chan struct{} {
	index := (atomic.LoadInt32(&this.curIndex) + this.interval - 1) % this.interval
	return this.timerC[index].Get()
}

func (this *TimeWheel) AfterFunc(f func()) {
	index := (atomic.LoadInt32(&this.curIndex) + this.interval - 1) % this.interval
	ml := &this.tasks[index]
	ml.Lock()
	ml.Elems.PushBack(f)
	ml.Unlock()
}

func (this *TimeWheel) Stop() {
	atomic.StoreInt32(&this.status, 1)
}

func (this *TimeWheel) onTicker() {
	tw := this
	curIndex := tw.curIndex
	this.timerC[curIndex].Notify()
	nextIndex := atomic.AddInt32(&tw.curIndex, 1) % tw.interval
	atomic.StoreInt32(&tw.curIndex, nextIndex)
	ml := &tw.tasks[curIndex]

	var elems list.List
	ml.Lock()
	elems = ml.Elems
	ml.Elems.Init()
	ml.Unlock()
	if elems.Len() > 1000 {
		go func(elems list.List) {
			e := elems.Front()
			if e != nil {
				for ; e != nil; e = e.Next() {
					e.Value.(func())()
				}
			}
		}(elems)
	} else {
		e := elems.Front()
		if e != nil {
			for ; e != nil; e = e.Next() {
				e.Value.(func())()
			}
		}
	}
}
