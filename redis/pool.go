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
		if this.Conn.Err() != nil {
			atomic.AddInt32(&this.pool.curActive, -1)
			return this.Conn.Close()
		}
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
		ca := atomic.LoadInt32(&this.curActive)
		if ca < this.maxActive {
			conn, err := this.callback()
			if err != nil {
				return nil, err
			}
			atomic.AddInt32(&this.curActive, 1)
			return &RedisConn{
				Conn: conn,
				pool: this,
			}, nil
		} else {
			fmt.Println("Error 0001 : too many active conn, maxActive=", this.maxActive)
			e := <-this.elems
			fmt.Println("return e")
			return e, nil
		}
	}
}

func (this *Pool) GetAsync() (*RedisConn, error) {
	if atomic.LoadInt32(&this.status) != 0 {
		return nil, fmt.Errorf("Error 0002 : this pool has been closed")
	}
	select {
	case e := <-this.elems:
		return e, nil
	default:
		ca := atomic.LoadInt32(&this.curActive)
		if ca < this.maxActive {
			conn, err := this.callback()
			if err != nil {
				return nil, err
			}
			atomic.AddInt32(&this.curActive, 1)
			return &RedisConn{
				Conn: conn,
				pool: this,
			}, nil
		} else {
			return nil, fmt.Errorf("Error 0001 : too many active conn, maxActive=%d", this.maxActive)
		}
	}
}

func (this *Pool) GetSync() (*RedisConn, error) {
	if atomic.LoadInt32(&this.status) != 0 {
		return nil, fmt.Errorf("Error 0002 : this pool has been closed")
	}
	return <-this.elems, nil
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
			default:
				break
			}
		}
	}
}
