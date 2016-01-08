package utils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type PoolElemInterface interface {
	Recycle() //回收
	Close()
	Err() error
	SetErr(error)
	Active()         //激活
	Heartbeat()      //心跳
	Timeout()        //设置超时，激活心跳
	IsTimeout() bool //是否超时
}

type PoolElem struct {
	Error  error
	Mux    sync.Mutex
	status int32
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
	callback    func(*Pool) (PoolElemInterface, error)
	elems       chan PoolElemInterface
	maxIdle     int32
	maxActive   int32
	curActive   int32
	timer       int32
	status      int32 //0正常
	elemsSize   int32
	timerStatus int32
}

func NewPool(callback func(*Pool) (PoolElemInterface, error), maxIdle, maxActive, timer int32) *Pool {
	pool := &Pool{
		callback:  callback,
		elems:     make(chan PoolElemInterface, maxActive),
		maxIdle:   maxIdle,
		maxActive: maxActive,
		timer:     timer,
	}
	if timer > 0 {
		go pool.timerEvent()
	}
	return pool
}

func (this *Pool) Close() {
	if this.IsClosed() {
		return
	}
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

func (this *Pool) IsClosed() bool {
	return atomic.LoadInt32(&this.status) != 0
}

func (this *Pool) Update(maxIdle, maxActive int32) {
	if this.IsClosed() {
		return
	}
	if maxIdle == this.maxIdle && maxActive == this.maxActive {
		return
	}
	this.maxIdle = maxIdle
	elems := this.elems
	this.elems = make(chan PoolElemInterface, maxActive)
	atomic.StoreInt32(&this.elemsSize, 0)
	flag := true
	for flag {
		select {
		case e := <-elems:
			select {
			case this.elems <- e:
				atomic.AddInt32(&this.elemsSize, 1)
			default:
				flag = false
			}
		default:
			flag = false
		}
	}
	atomic.StoreInt32(&this.maxActive, maxActive)
}

func (this *Pool) Put(elem PoolElemInterface) {
	if this.IsClosed() {
		elem.Close()
		return
	}

	if elem.Err() != nil {
		atomic.AddInt32(&this.curActive, -1)
		elem.Close()
		return
	}

	select {
	case this.elems <- elem:
		atomic.AddInt32(&this.elemsSize, 1)
		break
	default:
		atomic.AddInt32(&this.curActive, -1)
		elem.Close()
	}
}
func (this *Pool) Get() (PoolElemInterface, error) {
	if this.IsClosed() {
		return nil, fmt.Errorf("Error pool has been closed!")
	}
	var (
		elem PoolElemInterface
		err  error
	)
	for {
		elem, err = this.get()
		if err == nil && elem != nil && elem.Err() != nil {
			atomic.AddInt32(&this.curActive, -1)
			elem.Close()
			continue
		}
		break
	}
	return elem, err
}
func (this *Pool) get() (PoolElemInterface, error) {
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
		atomic.AddInt32(&this.elemsSize, -1)
	default:
		ca := atomic.LoadInt32(&this.curActive)
		if ca < this.maxActive {
			conn, err = this.callback(this)
			if err == nil {
				atomic.AddInt32(&this.curActive, 1)
			}
		} else {
			fmt.Println("Error 0001 : too many active conn, maxActive=", this.maxActive)
			select {
			case conn = <-this.elems:
				atomic.AddInt32(&this.elemsSize, -1)
				fmt.Println("return e")
			case <-time.After(time.Second * 2):
				return nil, fmt.Errorf("Error 0001 : too many active conn, maxActive=%d", this.maxActive)
			}
		}
	}
	if err == nil && conn != nil {
		conn.Active()
	}
	return conn, err
}

func (this *Pool) timerEvent() {
	timer := time.NewTicker(time.Second * time.Duration(this.timer))
	defer timer.Stop()
	for !this.IsClosed() {
		select {
		case <-timer.C:
			if atomic.LoadInt32(&this.elemsSize) > this.maxIdle {
				this.timerStatus++
				if this.timerStatus > 3 {
					select {
					case e := <-this.elems:
						atomic.AddInt32(&this.curActive, -1)
						atomic.AddInt32(&this.elemsSize, -1)
						e.Close()
					default:
						this.timerStatus = 0
					}
				} else {
					this.timerStatus = 0
				}
			}
			n := int(atomic.LoadInt32(&this.elemsSize)/30 + 1)
			flag := true
			for i := 0; i < n && flag; i++ {
				select {
				case e := <-this.elems:
					atomic.AddInt32(&this.elemsSize, -1)
					e.Heartbeat()
					e.Recycle()
				default:
					flag = false
					break
				}
			}
		}
	}
}
