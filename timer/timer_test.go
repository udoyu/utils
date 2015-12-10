package timer

import (
	"testing"
	"time"
)

func Test_TimeWheel(t *testing.T) {
	t.Log(time.Now().String())
	tw := NewTimeWheel(time.Second, 2)
	ch := make(chan bool, 1)
	tw.Add(func(){
		ch <- true
		t.Log(time.Now().String())
		})
 	<-ch
}

func Test_TimerHandler(t *testing.T) {
	timer := NewTimerHandler()
	ch2 := make(chan bool, 2)
	t.Log(time.Now().String())
	timer.Add(time.Second * 2, func() {
		ch2 <- true
		t.Log(time.Now().String())
	})
	timer.Add(time.Second * 3, func() {
		ch2 <- true
		t.Log(time.Now().String())
	})
	<-ch2
	<-ch2
}
