package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_TimeWheelAfterFunc(t *testing.T) {
	tw := NewTimeWheel(time.Second, 2)
	defer tw.Stop()
	i := int32(0)
	tw.AfterFunc(func() {
		atomic.AddInt32(&i, 1)
	})
	c := time.After(time.Second * 3)
	<-c
	if atomic.LoadInt32(&i) != 1 {
		t.Error(i)
	}
}

func Test_TimeWheelAfterFunc1(t *testing.T) {
	tw := NewTimeWheel(time.Second, 2)
	defer tw.Stop()
	i := int32(0)
	for j := 0; j < 10000; j++ {
		tw.AfterFunc(func() {
			atomic.AddInt32(&i, 1)
		})
	}
	c := time.After(time.Second * 2)
	<-c
	if atomic.LoadInt32(&i) != 10000 {
		t.Fatal(i)
	}
}

func Benchmark_TimerWheelAfterFunc(b *testing.B) {
	tw := NewTimeWheel(time.Second, 10)
	for i := 0; i < b.N; i++ {
		tw.AfterFunc(func() {})
	}
}

func Benchmark_MutexTicker(b *testing.B) {
	t := NewMutexTicker()
	for i:=0; i<b.N; i++ {
		t.Notify()
	}
}