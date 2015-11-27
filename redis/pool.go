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
	callback  func() (redis.Conn, error)
	elems     chan *RedisConn
	maxIdle   int32
	maxActive int32
	curActive int32
	status    int32 //1-closed
}

func NewPool(callback func() (redis.Conn, error), maxIdle, maxActive int32) *Pool {
	pool := &Pool{
		callback:  callback,
		elems:     make(chan *RedisConn, maxIdle),
		maxIdle:   maxIdle,
		maxActive: maxActive,
	}
	go pool.timerEvent()
	return pool
}

func (this *Pool) Update(maxIdle, maxActive int32) {
	atomic.StoreInt32(&this.maxActive, maxActive)
	if maxIdle == this.maxIdle {
		return
	}
	this.maxIdle = maxIdle
	elems := this.elems
	this.elems = make(chan *RedisConn, maxIdle)
	
	flag := true
	for flag {
		select {
		case e := <-elems:
			select {
				case this.elems <- e :
				default :
					flag = false
			}
		default:
			flag = false
		}
	}
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
	atomic.AddInt32(&this.curActive, -1)
	if elem.Conn.Err() != nil {
		elem.Conn.Close()
		return
	}

	select {
	case this.elems <- elem:
		break
	default:
		elem.Conn.Close()
	}

}

func (this *Pool) Get() (*RedisConn, error) {
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
		} else {
			fmt.Println("Error 0001 : too many active conn, maxActive=", this.maxActive)
			conn = <-this.elems
			fmt.Println("return e")
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
			e.Conn.Close()
		default:
			return
		}
	}
}

func (this *Pool) timerEvent() {
	timer := time.NewTicker(time.Second * 3)
	defer timer.Stop()
	for this.status == 0 {
		select {
		case <-timer.C:
			select {
			case e := <-this.elems:
				e.Do("PING")
				e.Close()
			default:
				break
			}
		}
	}
}
