package redis

import (
	"github.com/garyburd/redigo/redis"
	"reflect"
	"testing"
	"time"
)

var (
	pool      *Pool
	maxIdle   = int32(16)
	maxActive = int32(1024)
)

func TestNewPool(t *testing.T) {
	pool = NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "192.168.1.202:6379")
		if err != nil {
			t.Fatal(err)
		}
		return c, err
	},
		maxIdle,
		maxActive)
	conn, err := pool.Get()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
}

func TestGet(t *testing.T) {
	ch := make([]*RedisConn, 1024)
	var err error
	for i := 0; i < 1024; i++ {
		ch[i], err = pool.Get()
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 1024; i++ {
		pool.Put(ch[i])
	}
}

func TestDo(t *testing.T) {
	reply, err := pool.Do("SET", "test", "test")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reply.(string); !ok {
		t.Fatal("wrong reply type|type=", reflect.TypeOf(reply))
	} else if reply.(string) != "OK" {
		t.Fatal("reply wrong|reply=", reply)
	}
	reply, err = pool.Do("GET", "test")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := reply.([]uint8); !ok {
		t.Fatal("wrong reply type|type=", reflect.TypeOf(reply))
	} else if string(reply.([]uint8)) != "test" {
		t.Fatal("reply wrong|reply=", reply)
	}
	pool.Do("DEL", "test")
}

func TestTimerEvent(t *testing.T) {
	pool.Close()
	maxIdle = int32(3)
	maxActive = int32(8)
	pool = NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "192.168.1.202:6379")
		if err != nil {
			t.Fatal(err)
		}
		return c, err
	},
		maxIdle,
		maxActive)
	pool.SetLifeTime(0)
	elems := make([]*RedisConn, maxActive)
	var err error
	for i := 0; i < int(maxActive); i++ {
		elems[i], err = pool.Get()
		if err != nil {
			t.Fatal(err)
		}
	}
	if pool.curActive != maxActive {
		t.Fatal("size wrong|curActive=", pool.curActive, "|maxActive=", maxActive)
	}
	for i := 0; i < int(maxActive); i++ {
		elems[i].Close()
	}
	if pool.elemsSize != maxActive {
		t.Fatal("size wrong|elemsSize=", pool.elemsSize, "|maxActive=", maxActive)
	}
	time.AfterFunc(time.Second*3, func() {
		if pool.elemsSize != maxActive-3 || pool.curActive != maxActive-3 {
			t.Fatal("elemsSize=", pool.elemsSize, "|curActive=", pool.curActive)
		}
	})
	time.Sleep(time.Second * 4)
}

func TestWait(t *testing.T) {
	pool.Update(1,1)
	pool.SetWaitTime(1)
	{
	_,e := pool.Get()
	if e != nil {
		t.Error(e)
	}
	}
	{
	_, e := pool.Get()
	if e == nil {
		t.Error("failed")
	}
	}
}