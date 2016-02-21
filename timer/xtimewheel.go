package timer

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type TaskNode struct {
	activeTime time.Duration
	task       func()
}

type XMutexList struct {
	sync.Mutex
	Elems     []*list.List
	ElenIndex int
	c         chan struct{}
}

type XTimeWheel struct {
	ticker     *time.Ticker
	tasks      [][]XMutexList
	precisions []time.Duration
	intervals  []int64
	curIndexs  []int64
	bucket_cnt int
	status     int32
	offset     []int64
	tickets    []time.Duration
	pre_base   []int64
	now        time.Duration
}

func NewXTimeWheel(basetime time.Duration, intervals []int64) *XTimeWheel {
	tw := &XTimeWheel{}
	tw.bucket_cnt = len(intervals)
	tw.intervals = intervals

	tw.precisions = make([]time.Duration, tw.bucket_cnt)
	tw.pre_base = make([]int64, tw.bucket_cnt)
	tw.tickets = make([]time.Duration, tw.bucket_cnt)
	tw.precisions[0] = basetime
	for i := 0; i < tw.bucket_cnt; i++ {

		tw.precisions[i] = basetime
		tw.pre_base[i] = 1
		for j := 0; j < i; j++ {
			tw.precisions[i] *= time.Duration(tw.intervals[j])
			tw.pre_base[i] *= tw.intervals[j]
		}
		tw.tickets[i] = tw.precisions[i] * time.Duration(tw.intervals[i])
	}

	tw.curIndexs = make([]int64, tw.bucket_cnt)
	tw.offset = make([]int64, tw.bucket_cnt)

	tw.tasks = make([][]XMutexList, tw.bucket_cnt)
	for i := 0; i < tw.bucket_cnt; i++ {
		tw.tasks[i] = make([]XMutexList, tw.intervals[i])
	}
	tw.start()
	return tw
}
func (tw *XTimeWheel) onTimer(i int) {
	curIndex := tw.curIndexs[i]
	atomic.StoreInt64(&tw.curIndexs[i], (curIndex+1)%tw.intervals[i])

	ml := &tw.tasks[i][curIndex]

	var elems []*list.List
	var c chan struct{} = nil
	ml.Lock()
	c = ml.c
	elems = ml.Elems
	ml.Elems = nil
	ml.ElenIndex = 0
	ml.Unlock()
	if c != nil {
		close(c)
	}

	for _, v := range elems {
		go func(elems *list.List, tw *XTimeWheel, i int) {
			e := elems.Front()
			if e != nil {
				for ; e != nil; e = e.Next() {
					tn := e.Value.(*TaskNode)
					nextTime := tn.activeTime % tw.precisions[i]
					if nextTime == 0 ||
						i == 0 {
						tn.task()
					} else {
						tw.AfterFunc(nextTime, tn.task)
					}
				}
			}
		}(v, tw, i)
	}
}
func (this *XTimeWheel) start() {
	go func(tw *XTimeWheel) {
		tw.ticker = time.NewTicker(tw.precisions[0])
		defer tw.ticker.Stop()
		for atomic.LoadInt32(&tw.status) == 0 {
			select {
			case <-tw.ticker.C:
				for i := 0; i < this.bucket_cnt; i++ {
					if tw.UpdateOffset(i) == 0 {
						go tw.onTimer(i)
					} 
				}
			}
		}
	}(this)
}

func (this *XTimeWheel) UpdateOffset(index int) int64 {
	i := (this.offset[index] + 1) % int64(this.pre_base[index])
	atomic.StoreInt64(&this.offset[index], i)
	return i
}

func (this *XTimeWheel) After(d time.Duration) <-chan struct{} {
	var i = 0
	for i = 0; i < this.bucket_cnt-1; i++ {
		if d < this.precisions[i+1] {
			break
		}
	}
	interval := int64(d / this.precisions[i])
	if interval > this.intervals[i] {
		panic(fmt.Errorf("TimeWheel wrong after time, interval=%d and aftertime=%d",
			this.intervals[i]*int64(this.precisions[i]), d))
	} else if interval == 0 && i == 0 {
		return nil
	}

	index := (atomic.LoadInt64(&this.curIndexs[i]) + interval - 1) % this.intervals[i]
	ml := &this.tasks[i][index]
	var c chan struct{} = nil
	ml.Lock()
	if i == 0 {
		if ml.c == nil {
			ml.c = make(chan struct{})
		}
		c = ml.c
	} else {
		c = make(chan struct{})
		f := func() { close((c)) }
		if len(ml.Elems) == 0 || ml.Elems[ml.ElenIndex-1].Len() > 1000 {
			ml.Elems = append(ml.Elems, list.New())
			ml.ElenIndex++
		}
		d += time.Duration(this.offset[i]) * this.precisions[0]
		ml.Elems[ml.ElenIndex-1].PushBack(&TaskNode{activeTime: d, task: f})
	}
	ml.Unlock()
	return c
}

func (this *XTimeWheel) AfterFunc(d time.Duration, f func()) {
	var i = 0
	for i = 0; i < this.bucket_cnt-1; i++ {
		if d < this.precisions[i+1] {
			break
		}
	}
	d += time.Duration(atomic.LoadInt64(&this.offset[i])) * this.precisions[0]
	interval := int64(d / this.precisions[i])
	if interval > this.intervals[i] {
		panic(fmt.Errorf("TimeWheel wrong after time, interval=%d and aftertime=%d",
			this.intervals[i]*int64(this.precisions[i]), d))
	} else if interval == 0 && i == 0 {
		go f()
	}

	index := (atomic.LoadInt64(&this.curIndexs[i]) + interval - 1) % this.intervals[i]
	ml := &this.tasks[i][index]
	ml.Lock()
	if len(ml.Elems) == 0 || ml.Elems[ml.ElenIndex-1].Len() > 1000 {
		ml.Elems = append(ml.Elems, list.New())
		ml.ElenIndex++
	}

	ml.Elems[ml.ElenIndex-1].PushBack(&TaskNode{activeTime: d, task: f})
	ml.Unlock()
}

func (this *XTimeWheel) Stop() {
	atomic.StoreInt32(&this.status, 1)
}
