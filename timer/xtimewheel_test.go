package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_XTimeWheelAfterFunc(t *testing.T) {
	tw := NewXTimeWheel(time.Second, 2, 256)
	defer tw.Stop()
	i := int32(0)
	tw.AfterFunc(func() {
		atomic.AddInt32(&i, 1)
	})
	c := time.After(time.Second * 2)
	<-c
	if atomic.LoadInt32(&i) != 1 {
		t.Fatal(i)
	}
}

func Test_XTimeWheelAfterFunc1(t *testing.T) {
	tw := NewXTimeWheel(time.Second, 2, 256)
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