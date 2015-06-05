package simredis

import (
	"github.com/garyburd/redigo/redis"
	"sync"
)

type RedisProviderInfo struct {
	RedisAddrs   []string
	RedisDbNum   int
	CurRedisAddr string //当前连接的redis
}

type RedisProvider struct {
	RedisProviderInfo
	size     int
	pool     *redis.Pool
	index    int
	redisNum int
	lock     *sync.Mutex
}

//addresses[0]为使用地址，其他为备用地址
func NewRedisProvider(addresses []string, size, dbNum int) *RedisProvider {
	return &RedisProvider{
		RedisProviderInfo: RedisProviderInfo{
			RedisAddrs: addresses,
			RedisDbNum: dbNum,
		},
		size:     size,
		redisNum: len(addresses),
		lock:     new(sync.Mutex),
	}
}

func (this *RedisProvider) Pool() *redis.Pool {
	return this.pool
}

func (this *RedisProvider) Init() {
	this.Update()
}

func (this *RedisProvider) Update() {
	this.lock.Lock()
	if this.index >= this.redisNum {
		this.index = 0
	}
	i := this.index
	this.CurRedisAddr = this.RedisAddrs[i]
	this.index++
	this.lock.Unlock()

	this.pool = redis.NewPool(func() (redis.Conn, error) {
		println("curRedisInfo=", this.CurRedisAddr)
		c, e := redis.Dial("tcp", this.CurRedisAddr)
		if e != nil {
			return nil, e
		}
		_, e = c.Do("SELECT", this.RedisDbNum)
		if e != nil {
			c.Close()
			return nil, e
		}
		return c, e
	}, this.size)
}

func (this *RedisProvider) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	for i := 0; i < this.redisNum; i++ {
		if this.pool == nil {
			this.Update()
			if this.pool == nil {
				continue
			}
		}
		c := this.pool.Get()
		defer c.Close()
		reply, err = c.Do(commandName, args...)
		if err != nil {
			println(err.Error())
			this.Update()
		} else {
			break
		}
	}
	return reply, err
}

func (this *RedisProvider) Send(commands []Command) error {
	var err error
	for i := 0; i < this.redisNum; i++ {
		c := this.pool.Get()
		defer c.Close()
		for _, v := range commands {
			err = c.Send(v.CommandName, v.Args...)
			if err != nil {
				break
			}
		}

		if err != nil || c.Flush() != nil {
			this.Update()
			continue
		} else {
			break
		}
	}
	return err
}
