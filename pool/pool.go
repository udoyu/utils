package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type PoolElemInterface interface {
	Recyle() //回收
	Close()
	Err() error
	SetErr(error)
	Active()         //激活
	Heartbeat()      //心跳
	Timeout()        //设置超时，激活心跳
	IsTimeout() bool //是否超时
}

type PoolElem struct {
	Pool   *Pool
	Error  error
	Mux    sync.Mutex
	status int32
}

func (this *PoolElem) Recyle() {
	this.Pool.Put(this)
}

func (this *PoolElem) Close() {

}

func (this *PoolElem) Err() error {
	this.Mux.Lock()
	err := this.Error
	this.Mux.Unlock()
	return err
}

func (this *PoolElem) SetErr(err error) {
	this.Mux.Lock()
	this.Error = err
	this.Mux.Unlock()
}

func (this *PoolElem) Active() {
	atomic.StoreInt32(&this.status, 0)
}

func (this *PoolElem) Heartbeat() {
	//ping
}

func (this *PoolElem) Timeout() {
	atomic.AddInt32(&this.status, 1)
}

func (this *PoolElem) IsTimeout() bool {
	return atomic.LoadInt32(&this.status) > 1
}

type Pool struct {
	callback    func(*Pool) PoolElemInterface
	elems       chan PoolElemInterface
	activeElems chan PoolElemInterface
	maxIdle     int32
	maxActive   int32
	curActive   int32
	timer       int32
	status      int32 //0正常
}

func NewPool(callback func(*Pool) PoolElemInterface, maxIdle, maxActive, timer int32) *Pool {
	pool := &Pool{
		callback:    callback,
		elems:       make(chan PoolElemInterface, maxIdle),
		activeElems: make(chan PoolElemInterface, maxActive),
		maxIdle:     maxIdle,
		maxActive:   maxActive,
		timer:       timer,
	}
	if timer > 0 {
		go pool.timerEvent()
	}
	return pool
}

func (this *Pool) Put(elem PoolElemInterface) {
	if atomic.LoadInt32(&this.status) != 0 {
		elem.Close()
		return
	}
	atomic.AddInt32(&this.curActive, -1)
	if elem.Err() != nil {
		elem.Close()
		return
	}

	select {
	case this.elems <- elem:
		break
	default:
		select {
		case this.activeElems <- elem:
		default:
			elem.Close()
		}
	}
}

func (this *Pool) Get() (PoolElemInterface, error) {
	var (
		conn PoolElemInterface
		err  error
	)
	if atomic.LoadInt32(&this.status) != 0 {
		return conn, err
	}
	select {
	case e := <-this.elems:
		conn = e
	default:
		select {
		case e := <-this.activeElems:
			conn = e
		default:
			ca := atomic.LoadInt32(&this.curActive)
			if ca < this.maxActive {
				conn = this.callback(this)
			} else {
				fmt.Println("Error 0001 : too many active conn, maxActive=", this.maxActive)
				conn = <-this.elems
				fmt.Println("return e")
			}
		}

	}
	if conn != nil {
		atomic.AddInt32(&this.curActive, 1)
	}
	return conn, err
}

func (this *Pool) Close() {
	atomic.StoreInt32(&this.status, 1)
	for {
		select {
		case e := <-this.elems:
			e.Close()
		default:
			return
		}
	}
}

func (this *Pool) timerEvent() {
	timer := time.NewTicker(time.Second * time.Duration(this.timer))
	defer timer.Stop()
	for atomic.LoadInt32(&this.status) == 0 {
		select {
		case <-timer.C:
			select {
			case e := <-this.elems:
				if e.IsTimeout() {
					e.Heartbeat()
					e.Active()
				} else {
					e.Timeout()
				}
				e.Recyle()
			default:
				break
			}
			select {
			case e := <-this.activeElems:
				if e.IsTimeout() { //回收多余的空闲的链接
					e.Close()
				} else {
					//e.Do("PING")
					e.Timeout()
					e.Recyle()
				}
			}
		}
	}
}
