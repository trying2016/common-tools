package cache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/trying2016/common-tools/log"
	"github.com/trying2016/common-tools/utils"
)

type RedisConfig struct {
	Name     string
	PoolSize int
	Host     string
	Port     int
}

var errClientNull = errors.New("Redis client is null")

func init() {

}

type Cache struct {
	Name        string
	redisClient *redis.Client
}

func NewCache(client *redis.Client, name string) *Cache {
	cache := &Cache{
		redisClient: client,
		Name:        name,
	}
	return cache
}

func Strings(reply []interface{}, err error) ([]string, error) {
	var arr []string
	for _, v := range reply {
		arr = append(arr, utils.ToString(v))
	}
	return arr, nil
}

func (cache *Cache) Close() {
	if cache.redisClient != nil {
		_ = cache.redisClient.Close()
		cache.redisClient = nil
	}
}

func (cache *Cache) GetClient() (*redis.Client, error) {
	if cache.redisClient == nil {
		return nil, errClientNull
	} else {
		return cache.redisClient, nil
	}
}

func (cache *Cache) HSet(key string, field string, value string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.HSet(cache.Name+key, field, value).Err()
	if err != nil {
		log.Error("hset error: %v (%s, %s, %s, %s)", err, cache.Name, key, field, value)
	}
	return
}

func (cache *Cache) HGet(key string, field string) (ret string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return "", err
	}
	cmd := client.HGet(cache.Name+key, field)
	ret, err = cmd.Result()
	if err != nil {
		log.Error("HGet error: %v (%s, %s, %s)", err, cache.Name, key, field)
	}
	return
}

func (cache *Cache) HMGet(key string, args ...interface{}) (rets []string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.HMGet(cache.Name + key)
	err = cmd.Err()
	if err != nil {
		log.Error("HMGet error: %v (%s, %s)", err, cache.Name, key)
	}
	return Strings(cmd.Result())
}

func (cache *Cache) HGetAll(key string) (map[string]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.HGetAll(key)
	rets, err := cmd.Result()
	if err != nil {
		log.Error("HGetAll error: %v (%s, %s)", err, cache.Name, key)
		return nil, err
	}
	return rets, err
}

func (cache *Cache) Get(key string) (ret string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return "", err
	}
	cmd := client.Get(cache.Name + key)
	ret, err = cmd.Result()
	if err != nil {
		log.Error("Get error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}

func (cache *Cache) Set(key string, value string, expiration time.Duration) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return err
	}
	cmd := client.Set(cache.Name+key, value, expiration)
	err = cmd.Err()
	if err != nil {
		log.Error("set error: %v (%s, %s, %s)", err, cache.Name, key, value)
	}
	return
}

func (cache *Cache) LLen(key string) int {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return 0
	}
	cmd := client.LLen(cache.Name + key)
	ret, err := cmd.Result()
	if err != nil {
		log.Error("LLen error: %v (%s, %s)", err, cache.Name, key)
	}
	return int(ret)
}

//
func (cache *Cache) LRange(key string, start int, stop int) (rets []string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.LRange(cache.Name+key, int64(start), int64(stop))
	if cmd.Err() != nil {
		log.Error("LRange error: %v (%s, %s)", cmd.Err(), cache.Name, key)
	}
	return cmd.Result()
}

func (cache *Cache) RPush(key, value string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return err
	}
	cmd := client.RPush(cache.Name+key, value)
	err = cmd.Err()
	if err != nil {
		log.Error("RPush error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}

// 删除key
func (cache *Cache) Del(key string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.Del(cache.Name + key).Err()
	if err != nil {
		log.Error("set error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}
func (cache *Cache) Incrby(key string, value int64) (ret int64, err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return 0, err
	}
	cmd := client.IncrBy(cache.Name+key, value)
	ret, err = cmd.Result()
	if err != nil {
		log.Error("INCRBY error: %v (%s, %s, %v)", err, cache.Name, key, value)
	}
	return
}

func (cache *Cache) Expire(key string, expireTime time.Duration) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.Expire(cache.Name+key, expireTime).Err()
	if err != nil {
		log.Error("Expire error: %v (%s, %s, %d)", err, cache.Name, key, expireTime)
	}
	return
}

// 获取所有key
func (cache *Cache) Keys(key string) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.Keys(cache.Name + key)
	if cmd.Err() != nil {
		log.Error("Keys error: %v (%s, %s)", cmd.Err(), cache.Name, key)
	}
	return cmd.Result()
}