package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func Test_TimeWheel(t *testing.T) {
	tw := NewTimeWheel(time.Second, 2)
	defer tw.Stop()
	i := int32(0)
	tw.Add(func() {
		atomic.AddInt32(&i, 1)
	})
	c := time.After(time.Second * 2)
	<-c
	if atomic.LoadInt32(&i) != 1 {
		t.Fatal(i)
	}
}

func Test_TimeWheel1(t *testing.T) {
	tw := NewTimeWheel(time.Second, 2)
	defer tw.Stop()
	i := int32(0)
	for j := 0; j < 10000; j++ {
		tw.Add(func() {
			atomic.AddInt32(&i, 1)
		})
	}
	c := time.After(time.Second * 2)
	<-c
	if atomic.LoadInt32(&i) != 10000 {
		t.Fatal(i)
	}
}

func Test_TimerHandler(t *testing.T) {
	timer := NewTimerHandler()
	defer timer.Stop()
	{
		i := int32(0)
		timer.Add(time.Second*2, func() {
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
		timer.Add(time.Second*3, func() {
			atomic.AddInt32(&i, 1)
		})
		c := time.After(time.Second * 3)
		<-c
		if atomic.LoadInt32(&i) != 1 {
			t.Fatal(i)
		}
	}
}
