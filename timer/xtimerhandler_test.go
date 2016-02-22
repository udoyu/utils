package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_XTimerHandlerAfterFunc(t *testing.T) {
	timer := NewXTimerHandler(10)
	defer timer.Stop()
	{
		i := int32(0)
		ch := make(chan bool)
		timer.AfterFunc(time.Second*2, func() {
			atomic.AddInt32(&i, 1)
			ch <- true
		})
		<-ch
	}
	//	{
	//		i := int32(0)
	//		timer.AfterFunc(time.Second*3, func() {
	//			atomic.AddInt32(&i, 1)
	//		})
	//		c := time.After(time.Second * 3)
	//		<-c
	//		if atomic.LoadInt32(&i) != 1 {
	//			t.Fatal(i)
	//		}
	//	}
}

func Benchmark_XTimerHandlerAfterFunc(b *testing.B) {
	timer := NewXTimerHandler(256)
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.AfterFunc(time.Second, func() {})
	}
}

func Benchmark_XTimerHandlerAfter(b *testing.B) {
	timer := NewXTimerHandler(256)
	defer timer.Stop()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.After(time.Second)
	}
}
