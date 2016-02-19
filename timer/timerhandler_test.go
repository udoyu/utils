package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_TimerHandlerAfterFunc(t *testing.T) {
	timer := NewTimerHandler()
	defer timer.Stop()
	{
		i := int32(0)
		timer.AfterFunc(time.Second*2, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 3)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Error(i)
		}
	}
	{
		i := int32(0)
		timer.AfterFunc(time.Second*3, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 3)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Error(i)
		}
	}
}


func Test_TimerHandlerAfter(t *testing.T) {
	timer := NewTimerHandler()
	defer timer.Stop()
	select {
		case <-timer.After(time.Second):
		case <-time.After(time.Millisecond * 1001):
			t.Error("failed")
	}
}

func Benchmark_TimerHandlerAfterFunc(b *testing.B) {
	timer := NewTimerHandler()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.AfterFunc(time.Second, func() {})
	}
}

func Benchmark_TimerHandlerAfter(b *testing.B) {
	timer := NewTimerHandler()
	defer timer.Stop()
	for i:=0; i<b.N; i++ {
		timer.After(time.Second)
	}
}