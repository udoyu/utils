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

type TimeWheel struct {
	ticket   *time.Ticker
	tasks    []MutexList
	interval int32
	curIndex int32
	status   int32
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
	go func(tw *TimeWheel) {
		defer tw.ticket.Stop()
		for atomic.LoadInt32(&tw.status) == 0 {
			select {
			case <-tw.ticket.C:
				curIndex := tw.curIndex
				nextIndex := atomic.AddInt32(&tw.curIndex, 1) % tw.interval
				atomic.StoreInt32(&tw.curIndex, nextIndex)
				ml := &tw.tasks[curIndex]
				var elems list.List
				ml.Lock()
				elems = ml.Elems
				ml.Elems.Init()
				ml.Unlock()
				for e := elems.Front(); e != nil; e = e.Next() {
					e.Value.(func())()
				}
			}
		}
	}(tw)
	return tw
}

func (this *TimeWheel) Add(task func()) {
	index := (atomic.LoadInt32(&this.curIndex) + this.interval - 1) % this.interval
	ml := &this.tasks[index]
	ml.Lock()
	ml.Elems.PushBack(task)
	ml.Unlock()
}

func (this *TimeWheel) Stop() {
	atomic.StoreInt32(&this.status, 1)
}
