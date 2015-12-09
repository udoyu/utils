package redis

import "github.com/garyburd/redigo/redis"
import "os"
import "fmt"

type RedisHandle struct {
	Addr string
	Pool *Pool
}

// XXX: add some password protection
func NewRedisHandle(addr string, max_idle, max_active int, debug bool) *RedisHandle {
	if debug {
		fmt.Println("[RedisHandle] Opening New Handle For Pid:", os.Getpid())
	}
	return &RedisHandle{
		Addr: addr,
		Pool: NewPool(func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
			int32(max_idle),
			int32(max_active)),
	}
}

func (self *RedisHandle) GetRedisConn() (*RedisConn, error) {
	return self.Pool.Get()
}

func (self *RedisHandle) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	rc, err := self.GetRedisConn()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return rc.Do(commandName, args...)
}

// XXX: is _not_ calling defer rc.Close()
//      so do it yourself later
func (self *RedisHandle) Send(cmd string, args ...interface{}) (err error) {
	rc, err := self.GetRedisConn()
	if err != nil {
		return err
	}
	return rc.Send(cmd, args...)
}
