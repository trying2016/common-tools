package cache

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/trying2016/common-tools/log"
	"github.com/trying2016/common-tools/utils"
)

type RedisConfig struct {
	Name     string
	PoolSize int
	Host     string
	Port     int
}

func init() {

}

type Cache struct {
	Name        string
	RedisClient *redis.Pool
	cfg         *RedisConfig
}

func NewCache(cfg *RedisConfig) *Cache {
	utils.NewHttpClient()
	cache := &Cache{}
	cache.cfg = cfg
	cache.Name = cfg.Name
	cache.Init()
	return cache
}

func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func (cache *Cache) Init() {
	cache.RedisClient = &redis.Pool{
		MaxIdle:     20,
		Wait:        true,
		MaxActive:   cache.cfg.PoolSize, // max pool size
		IdleTimeout: 30 * time.Second,   //timeout
		Dial: func() (redis.Conn, error) {
			dns := fmt.Sprintf("%s:%d", cache.cfg.Host, cache.cfg.Port)
			c, err := redis.Dial("tcp", dns)
			if err != nil {
				return nil, err
			}
			return c, nil
		}, TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (cache *Cache) Close() {
	if cache.RedisClient != nil {
		cache.RedisClient.Close()
		cache.RedisClient = nil
	}
}

func (cache *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	rd := cache.RedisClient.Get()
	defer rd.Close()
	reply, err = rd.Do(commandName, args...)
	return
}

func (cache *Cache) HSet(key string, field string, value string) (err error) {
	_, err = cache.do("HSET", cache.Name+key, field, value)
	if err != nil {
		log.Error("hset error: %s (%s, %s, %s, %s)", err.Error(), cache.Name, key, field, value)
	}
	return
}

func (cache *Cache) HGet(key string, field string) (ret string, err error) {
	ret, err = redis.String(cache.do("HGET", cache.Name+key, field))
	return
}

func (cache *Cache) HMGet(key string, args ...interface{}) (rets []string, err error) {
	rets, err = redis.Strings(cache.do("HMGET", cache.Name+key, args))
	return
}

func (cache *Cache) HGetAll(key string) (map[string]string, error) {
	arr, err := redis.Strings(cache.do("hgetall", cache.Name+key))
	if err != nil {
		log.Error("keys error: %s (%s, %s)", err.Error(), cache.Name, key)
		return nil, err
	} else {
		rets := make(map[string]string)
		for i := 0; i < len(arr); i += 2 {
			rets[arr[i]] = arr[i+1]
		}
		return rets, nil
	}
}

func (cache *Cache) Get(key string) (ret string, err error) {
	ret, err = redis.String(cache.do("GET", cache.Name+key))
	return
}

func (cache *Cache) Set(key string, value string) (err error) {
	_, err = cache.do("SET", cache.Name+key, value)
	if err != nil {
		log.Error("set error: %s (%s, %s, %s)", err.Error(), cache.Name, key, value)
	}
	return
}

func (cache *Cache) LLen(key string) int {
	ret, err := redis.Int(cache.do("LLEN", cache.Name+key))
	if err != nil {
		log.Error("LLen error: %s (%s, %s)", err.Error(), cache.Name, key)
	}
	return ret
}

//
func (cache *Cache) LRange(key string, start int, stop int) (rets []string, err error) {
	rets, err = Strings(cache.do("LRANGE", cache.Name+key, start, stop))
	if err != nil {
		log.Error("LRange error: %s (%s, %s)", err.Error(), cache.Name, key)
	}
	return
}

func (cache *Cache) RPush(key, value string) (err error) {
	_, err = cache.do("RPUSH", cache.Name+key, value)
	if err != nil {
		log.Error("RPush error: %s (%s, %s)", err.Error(), cache.Name, key)
	}
	return
}

// 删除key
func (cache *Cache) Del(key string) (err error) {
	_, err = cache.do("del", cache.Name+key)
	if err != nil {
		log.Error("set error: %s (%s, %s)", err.Error(), cache.Name, key)
	}
	return
}
func (cache *Cache) Incrby(key string, value int) (ret int, err error) {
	ret, err = redis.Int(cache.do("INCRBY", cache.Name+key, value))
	if err != nil {
		log.Error("INCRBY error: %s (%s, %s, %d)", err.Error(), cache.Name, key, value)
	}
	return
}

func (cache *Cache) Expire(key string, time int) (err error) {
	_, err = cache.do("EXPIRE", cache.Name+key, time)
	if err != nil {
		log.Error("Expire error: %s (%s, %s, %d)", err.Error(), cache.Name, key, time)
	}
	return
}

// 获取所有key
func (cache* Cache) Keys(key string) []string{
	arr, err := redis.Strings(cache.do("keys", cache.Name+key))
	if err != nil {
		return nil
	}
	return arr
}