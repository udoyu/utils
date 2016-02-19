package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_AfterFunc(t *testing.T) {
	{
		i := int32(0)
		AfterFunc(time.Second*2, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 2)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Fatal(i)
		}
	}
	{
		i := int32(0)
		AfterFunc(time.Second*3, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 3)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Fatal(i)
		}
	}
}

func Test_After(t *testing.T) {
	select {
		case <-After(time.Second):
		case <-time.After(time.Millisecond * 1001):
			t.Error("failed")
	}
}

func Benchmark_AfterFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AfterFunc(time.Second, func() {})
	}
}

func Benchmark_After(b *testing.B) {
	for i:=0; i<b.N; i++ {
		After(time.Second)
	}
}
