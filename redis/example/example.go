package main

import (
	"fmt"
	redis "github.com/garyburd/redigo/redis"
	simredis "github.com/udoyu/utils/redis"
	"reflect"
)

const (
	REDIS_REPLY_STRING  = 1
	REDIS_REPLY_ARRAY   = 2
	REDIS_REPLY_INTEGER = 3
	REDIS_REPLY_NIL     = 4
	REDIS_REPLY_STATUS  = 5
	REDIS_REPLY_ERROR   = 6
)

type RedisReply struct {
	Type     int           /* REDIS_REPLY_* */
	Integer  int64         /* The integer when type is REDIS_REPLY_INTEGER */
	Len      int           /* Length of string */
	Str      string        /* Used for both REDIS_REPLY_ERROR and REDIS_REPLY_STRING */
	Elements int           /* number of elements, for REDIS_REPLY_ARRAY */
	Element  []*RedisReply /* elements vector for REDIS_REPLY_ARRAY */
}

func NewRedisReply(re interface{}, err error) *RedisReply {
	reply := &RedisReply{}
	if err != nil {
		reply.Type = REDIS_REPLY_ERROR
		reply.Str = err.Error()
		reply.Len = len(reply.Str)
		return reply
	}
	if re == nil {
		reply.Type = REDIS_REPLY_NIL
		return reply
	}
	switch re.(type) {
	case []uint8:
		reply.Type = REDIS_REPLY_STRING
		reply.Str = string(re.([]uint8))
		reply.Len = len(reply.Str)
	case []interface{}:
		reply.Type = REDIS_REPLY_ARRAY
		reply.Elements = len(re.([]interface{}))
		replys := make([]*RedisReply, reply.Elements)
		for i, r := range re.([]interface{}) {
			replys[i] = NewRedisReply(r, nil)
		}
		reply.Element = replys
	case int64:
		reply.Type = REDIS_REPLY_INTEGER
		reply.Integer = re.(int64)
	}
	return reply
}

type RedisInterface interface {
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
}

type Redis struct {
	value RedisInterface
}

func (this Redis) Do(commandName string, args ...interface{}) *RedisReply {
	re, err := this.value.Do(commandName, args...)
	fmt.Println(re, err)
	fmt.Println(reflect.TypeOf(re))
	return NewRedisReply(re, err)
}

func TestPool() {
	ri := simredis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		fmt.Println("Dial ...")
		if err != nil {
			fmt.Println(err)
		}
		return c, err
	},
		8,
		8)
	reply := Redis{ri}.Do("SET", "test", "hello")
	fmt.Println(reply)
	reply = Redis{ri}.Do("GET", "test")
	fmt.Println(*reply)
}

func TestCluster() {
	redishosts := []string{
		"127.0.0.1:6380",
		"127.0.0.1:6381",
		"127.0.0.1:6382",
		"127.0.0.1:6383",
		"127.0.0.1:6384",
		"127.0.0.1:6385",
	}
	ri := simredis.NewRedisCluster(redishosts, 8, 64, true)
	reply := Redis{&ri}.Do("SET", "test", "hello")
	fmt.Println(reply)
	reply = Redis{&ri}.Do("GET", "test")
	fmt.Println(*reply)
}

func main() {
	TestPool()
	//TestCluster()
}
