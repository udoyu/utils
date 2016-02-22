package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_TimerHandlerAfterFunc(t *testing.T) {
	timer := NewTimer()
	defer timer.Stop()
	{
		i := int32(0)
		ch := make(chan bool)
		timer.AfterFunc(time.Second*2, func() {
			atomic.AddInt32(&i, 1)
			ch <- true
		})
		<-ch
		//		<-time.After(time.Millisecond * 2100)
		//		if atomic.LoadInt32(&i) != 1 {
		//			t.Error(atomic.LoadInt32(&i))
		//		}
	}
}

func Test_TimerHandlerAfter(t *testing.T) {
	timer := NewTimer()
	defer timer.Stop()
	select {
	case <-timer.After(time.Second):
	case <-time.After(time.Millisecond * 1100):
		t.Error("failed")
	}
}

func Benchmark_TimerHandlerAfterFunc(b *testing.B) {
	timer := NewTimer()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.AfterFunc(time.Second, func() {})
	}
}

func Benchmark_TimerHandlerAfter(b *testing.B) {
	timer := NewTimer()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		timer.After(time.Second)
	}
}

func Benchmark_TimerHandlerAfterFuncAdd(b *testing.B) {
	timer := NewTimer()
	defer timer.Stop()
	bt := globalBaseTime * 2000
	for i := 0; i < b.N; i++ {
		timer.AfterFunc(bt, func() {})
	}
}

func Benchmark_TimerHandleron(b *testing.B) {
	timer := NewTimer()
	defer timer.Stop()
	for i := 0; i < b.N; i++ {
		//_ = bt%time.Duration(BASE_PER_BUCKET[1])
		for i := 0; i < 5; i++ {
			if timer.timerWheels.UpdateOffset(i) == 0 {
				timer.timerWheels.onTimer(i)
			} else {
				break
			}
		}
	}
}
