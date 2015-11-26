package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sync/atomic"
)

type RedisConn struct {
	redis.Conn
	pool      *Pool
	activeNum int32 //0，1,正常，2超时
}

func (this *RedisConn) Close() error {
	if this.Conn.Err() != nil {
		return this.Conn.Close()
	} else if this.pool != nil {
		this.pool.Put(this)
	} else {
		return this.Conn.Close()
	}
	return nil
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
	return &Pool{
		callback:  callback,
		elems:     make(chan *RedisConn, maxIdle),
		maxIdle:   maxIdle,
		maxActive: maxActive,
	}
}

func (this *Pool) Put(elem *RedisConn) {
	if atomic.LoadInt32(&this.status) != 0 {
		elem.Conn.Close()
	}
	select {
	case this.elems <- elem:
		break
	default:
		atomic.AddInt32(&this.curActive, -1)
		elem.Conn.Close()
	}
}

func (this *Pool) Get() (*RedisConn, error) {
	if atomic.LoadInt32(&this.status) != 0 {
		return nil, fmt.Errorf("Error 0002 : this pool has been closed")
	}
	select {
	case e := <-this.elems:
		return e, nil
	default:
		ca := atomic.AddInt32(&this.curActive, 1)
		if ca < this.maxActive {
			conn, err := this.callback()
			if err != nil {
				return nil, err
			}
			return &RedisConn{
				Conn: conn,
				pool: this,
			}, nil
		} else {
			return nil, fmt.Errorf("Error 0001 : too many active conn, maxActive=%d", this.maxActive)
		}
	}
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
