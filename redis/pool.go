package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync/atomic"
	"time"
)

type RedisConn struct {
	redis.Conn
	pool *Pool
	//	status int32
}

func (this *RedisConn) Close() error {
	if this.pool != nil {
		this.pool.Put(this)
	} else {
		return this.Conn.Close()
	}
	return nil
}

func (this RedisConn) Command(commandName string, args ...interface{}) *RedisReply {
	reply, err := this.Do(commandName, args)
	return NewRedisReply(reply, err)
}

type Pool struct {
	callback    func() (redis.Conn, error)
	elems       chan *RedisConn
	maxIdle     int32
	maxActive   int32
	curActive   int32
	elemsSize   int32
	status      int32 //1-closed
	timerStatus int32
}

func NewPool(callback func() (redis.Conn, error), maxIdle, maxActive int32) *Pool {
	pool := &Pool{
		callback:  callback,
		elems:     make(chan *RedisConn, maxActive),
		maxIdle:   maxIdle,
		maxActive: maxActive,
	}
	go pool.timerEvent()
	return pool
}

func (this *Pool) Update(maxIdle, maxActive int32) {

	if maxIdle == this.maxIdle && maxActive == this.maxActive {
		return
	}
	this.maxIdle = maxIdle
	elems := this.elems
	this.elems = make(chan *RedisConn, maxActive)
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

func (this *Pool) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	c, err := this.Get()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.Do(commandName, args...)
}

func (this *Pool) Put(elem *RedisConn) {
	if atomic.LoadInt32(&this.status) != 0 {
		elem.Conn.Close()
		return
	}

	if elem.Conn.Err() != nil {
		atomic.AddInt32(&this.curActive, -1)
		elem.Conn.Close()
		return
	}

	select {
	case this.elems <- elem:
		atomic.AddInt32(&this.elemsSize, 1)
		break
	default:
		elem.Conn.Close()
		atomic.AddInt32(&this.curActive, -1)
	}
}

func (this *Pool) Get() (*RedisConn, error) {
	var (
		elem *RedisConn
		err  error
	)
	for {
		elem, err = this.get()
		if elem != nil && elem.Err() != nil {
			atomic.AddInt32(&this.curActive, -1)
			elem.Conn.Close()
			continue
		}
		break
	}
	return elem, err
}

func (this *Pool) get() (*RedisConn, error) {
	if atomic.LoadInt32(&this.status) != 0 {
		return nil, fmt.Errorf("Error 0002 : this pool has been closed")
	}
	var (
		conn *RedisConn
		err  error
	)
	select {
	case e := <-this.elems:
		conn = e
		atomic.AddInt32(&this.elemsSize, -1)
	default:
		ca := atomic.LoadInt32(&this.curActive)
		if ca < this.maxActive {
			var c redis.Conn
			c, err = this.callback()
			if err != nil {
				break
			}

			conn = &RedisConn{
				Conn: c,
				pool: this,
			}
			atomic.AddInt32(&this.curActive, 1)
		} else {
			fmt.Println("Error 0001 : too many active conn, maxActive=", this.maxActive)
			select {
			case conn = <-this.elems:
				atomic.AddInt32(&this.elemsSize, -1)
				fmt.Println("return e")
			case <-time.After(time.Second * 3):
				err = fmt.Errorf("Error 0003 : RedisPool Get timeout")
			}
		}

	}
	//	if conn != nil {
	//		atomic.StoreInt32(&conn.status, 0)
	//	}
	return conn, err
}

func (this *Pool) Close() {
	atomic.StoreInt32(&this.status, 1)
	for {
		select {
		case e := <-this.elems:
			e.Conn.Close()
		default:
			return
		}
	}
}

func (this *Pool) timerEvent() {
	timer := time.NewTicker(time.Second * 3)
	defer timer.Stop()
	for atomic.LoadInt32(&this.status) == 0 {

		select {
		case <-timer.C:
			if atomic.LoadInt32(&this.elemsSize) > this.maxIdle {
				this.timerStatus++
				if this.timerStatus > 3 {
					select {
					case e := <-this.elems:
						atomic.AddInt32(&this.curActive, -1)
						atomic.AddInt32(&this.elemsSize, -1)
						e.Conn.Close()
					default:
						this.timerStatus = 0
					}
				} else {
					this.timerStatus = 0
				}
			}
			select {
			case e := <-this.elems:
				atomic.AddInt32(&this.elemsSize, -1)
				e.Do("PING")
				e.Close()
			default:
				break
			}
		}
	}
}
