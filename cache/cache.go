package cache

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
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

func (cache *Cache) HSet(key string, field string, value interface{}) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.HSet(cache.Name+key, field, value).Err()
	if err != nil {
		//log.Error("hset error: %v (%s, %s, %s, %s)", err, cache.Name, key, field, value)
	}
	return
}
func (cache *Cache) HMSet(key string, fields map[string]interface{}) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.HMSet(cache.Name+key, fields).Err()
	return err
}
func (cache *Cache) HGet(key string, field string) (ret string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return "", err
	}
	cmd := client.HGet(cache.Name+key, field)
	ret, err = cmd.Result()
	if err != nil {
		//log.Error("HGet error: %v (%s, %s, %s)", err, cache.Name, key, field)
	}
	return
}

func (cache *Cache) HMGet(key string, args ...string) (rets []string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.HMGet(cache.Name+key, args...)
	err = cmd.Err()
	if err != nil {
		//log.Error("HMGet error: %v (%s, %s)", err, cache.Name, key)
	}
	return Strings(cmd.Result())
}

func (cache *Cache) HGetAll(key string) (map[string]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.HGetAll(cache.Name + key)
	rets, err := cmd.Result()
	if err != nil {
		//log.Error("HGetAll error: %v (%s, %s)", err, cache.Name, key)
		return nil, err
	}
	return rets, err
}

func (cache *Cache) Get(key string) (ret string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return "", err
	}
	cmd := client.Get(cache.Name + key)
	ret, err = cmd.Result()
	if err != nil {
		//log.Error("Get error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}

func (cache *Cache) Set(key string, value interface{}, expiration time.Duration) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	cmd := client.Set(cache.Name+key, value, expiration)
	err = cmd.Err()
	if err != nil {
		//log.Error("set error: %v (%s, %s, %s)", err, cache.Name, key, value)
	}
	return
}

func (cache *Cache) LLen(key string) int {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return 0
	}
	cmd := client.LLen(cache.Name + key)
	ret, err := cmd.Result()
	if err != nil {
		//log.Error("LLen error: %v (%s, %s)", err, cache.Name, key)
	}
	return int(ret)
}

func (cache *Cache) LRange(key string, start int, stop int) (rets []string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.LRange(cache.Name+key, int64(start), int64(stop))
	if cmd.Err() != nil {
		//log.Error("LRange error: %v (%s, %s)", cmd.Err(), cache.Name, key)
	}
	return cmd.Result()
}

func (cache *Cache) RPush(key, value string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	cmd := client.RPush(cache.Name+key, value)
	err = cmd.Err()
	if err != nil {
		//log.Error("RPush error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}

func (cache *Cache) LPush(key, value string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	cmd := client.LPush(cache.Name+key, value)
	err = cmd.Err()
	if err != nil {
		//log.Error("RPush error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}

func (cache *Cache) LPop(key string) (str string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		return "", err
	}
	cmd := client.LPop(cache.Name + key)
	return cmd.Result()
}

func (cache *Cache) RPop(key string) (str string, err error) {
	client, err := cache.GetClient()
	if err != nil {
		return "", err
	}
	cmd := client.RPop(cache.Name + key)
	return cmd.Result()
}

// 删除key
func (cache *Cache) Del(key string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.Del(cache.Name + key).Err()
	if err != nil {
		//log.Error("set error: %v (%s, %s)", err, cache.Name, key)
	}
	return
}
func (cache *Cache) Incrby(key string, value int64) (ret int64, err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return 0, err
	}
	cmd := client.IncrBy(cache.Name+key, value)
	ret, err = cmd.Result()
	if err != nil {
		//log.Error("INCRBY error: %v (%s, %s, %v)", err, cache.Name, key, value)
	}
	return
}

func (cache *Cache) Expire(key string, expireTime time.Duration) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return err
	}
	err = client.Expire(cache.Name+key, expireTime).Err()
	if err != nil {
		//log.Error("Expire error: %v (%s, %s, %d)", err, cache.Name, key, expireTime)
	}
	return
}

// 获取所有key
func (cache *Cache) Keys(key string) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.Keys(cache.Name + key)
	if cmd.Err() != nil {
		//log.Error("Keys error: %v (%s, %s)", cmd.Err(), cache.Name, key)
	}
	return cmd.Result()
}

func (cache *Cache) Rename(key, newKey string) (err error) {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	cmd := client.Rename(cache.Name+key, cache.Name+newKey)
	return cmd.Err()
}

// 获取所有key
func (cache *Cache) HKeys(key string) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return nil, err
	}
	cmd := client.HKeys(cache.Name + key)
	if cmd.Err() != nil {
		//log.Error("Keys error: %v (%s, %s)", cmd.Err(), cache.Name, key)
	}
	return cmd.Result()
}
func (cache *Cache) HLen(key string) int {
	client, err := cache.GetClient()
	if err != nil {
		//log.Error("GetClient error: %v ", err)
		return 0
	}
	cmd := client.HLen(cache.Name + key)
	ret, err := cmd.Result()
	if err != nil {
		//log.Error("LLen error: %v (%s, %s)", err, cache.Name, key)
	}
	return int(ret)
}

func (cache *Cache) Subscribe(fn func(message *redis.Message, err error), channel string) error {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	ch := client.Subscribe(cache.Name + channel)
	go func() {
		for {
			msg, err := ch.ReceiveMessage()
			fn(msg, err)
		}
	}()

	return nil
}

// publish
func (cache *Cache) Publish(channel string, message interface{}) error {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	cmd := client.Publish(cache.Name+channel, message)
	return cmd.Err()
}

func (cache *Cache) HDel(key string, fields ...string) error {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	cmd := client.HDel(cache.Name+key, fields...)
	return cmd.Err()
}

// ZAdd 添加有序集合
func (cache *Cache) ZAdd(key string, score float64, member string) error {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	cmd := client.ZAdd(cache.Name+key, redis.Z{Score: score, Member: member})
	return cmd.Err()
}

// ZRange 获取有序集合
func (cache *Cache) ZRange(key string, start, stop int64) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.ZRange(cache.Name+key, start, stop)
	return cmd.Result()
}

// ZRevRange 获取有序集合
func (cache *Cache) ZRevRange(key string, start, stop int64) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.ZRevRange(cache.Name+key, start, stop)
	return cmd.Result()
}

// ZRem 移除有序集合
func (cache *Cache) ZRem(key string, members ...interface{}) error {
	client, err := cache.GetClient()
	if err != nil {
		return err
	}
	cmd := client.ZRem(cache.Name+key, members...)
	return cmd.Err()
}

// ZCard 获取有序集合的长度
func (cache *Cache) ZCard(key string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZCard(cache.Name + key)
	return cmd.Result()
}

// ZScore 获取有序集合的分数
func (cache *Cache) ZScore(key, member string) (float64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZScore(cache.Name+key, member)
	return cmd.Result()
}

// ZRank 获取有序集合的排名
func (cache *Cache) ZRank(key, member string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZRank(cache.Name+key, member)
	return cmd.Result()
}

// ZRevRank 获取有序集合的排名
func (cache *Cache) ZRevRank(key, member string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZRevRank(cache.Name+key, member)
	return cmd.Result()
}

// ZRangeByScore 获取有序集合的排名
func (cache *Cache) ZRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.ZRangeByScore(cache.Name+key, opt)
	return cmd.Result()
}

// ZRevRangeByScore 获取有序集合的排名
func (cache *Cache) ZRevRangeByScore(key string, opt redis.ZRangeBy) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.ZRevRangeByScore(cache.Name+key, opt)
	return cmd.Result()
}

// ZIncrBy 增加有序集合的分数
func (cache *Cache) ZIncrBy(key string, increment float64, member string) (float64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZIncrBy(cache.Name+key, increment, member)
	return cmd.Result()
}

// ZRemRangeByRank 移除有序集合的排名
func (cache *Cache) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZRemRangeByRank(cache.Name+key, start, stop)
	return cmd.Result()
}

// ZRemRangeByScore 移除有序集合的分数
func (cache *Cache) ZRemRangeByScore(key string, min, max string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.ZRemRangeByScore(cache.Name+key, min, max)
	return cmd.Result()
}

// ZUnionStore 合并有序集合
func (cache *Cache) ZUnionStore(dest string, store redis.ZStore, keys ...string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}

	fullKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		fullKeys = append(fullKeys, cache.Name+key)
	}

	cmd := client.ZUnionStore(cache.Name+dest, store, fullKeys...)
	return cmd.Result()
}

// ZInterStore 交集有序集合
func (cache *Cache) ZInterStore(dest string, store redis.ZStore, keys ...string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}

	fullKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		fullKeys = append(fullKeys, cache.Name+key)
	}

	cmd := client.ZInterStore(cache.Name+dest, store, fullKeys...)
	return cmd.Result()
}

// SAdd 添加集合
func (cache *Cache) SAdd(key string, members ...interface{}) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.SAdd(cache.Name+key, members...)
	return cmd.Result()
}

// SCard 获取集合的长度
func (cache *Cache) SCard(key string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.SCard(cache.Name + key)
	return cmd.Result()
}

// SDiffStore 差集集合
func (cache *Cache) SDiffStore(dest string, keys ...string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}

	fullKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		fullKeys = append(fullKeys, cache.Name+key)
	}

	cmd := client.SDiffStore(cache.Name+dest, fullKeys...)
	return cmd.Result()
}

// SInterStore 交集集合
func (cache *Cache) SInterStore(dest string, keys ...string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}

	fullKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		fullKeys = append(fullKeys, cache.Name+key)
	}

	cmd := client.SInterStore(cache.Name+dest, fullKeys...)
	return cmd.Result()
}

// SUnionStore 并集集合
func (cache *Cache) SUnionStore(dest string, keys ...string) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}

	fullKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		fullKeys = append(fullKeys, cache.Name+key)
	}

	cmd := client.SUnionStore(cache.Name+dest, fullKeys...)
	return cmd.Result()
}

// SMembers 获取集合的所有成员
func (cache *Cache) SMembers(key string) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.SMembers(cache.Name + key)
	return cmd.Result()
}

// SIsMember 判断成员是否在集合中
func (cache *Cache) SIsMember(key string, member string) (bool, error) {
	client, err := cache.GetClient()
	if err != nil {
		return false, err
	}
	cmd := client.SIsMember(cache.Name+key, member)
	return cmd.Result()
}

// SRem 移除集合的成员
func (cache *Cache) SRem(key string, members ...interface{}) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.SRem(cache.Name+key, members...)
	return cmd.Result()
}

// SMove 移动集合的成员
func (cache *Cache) SMove(src, dest string, member interface{}) (bool, error) {
	client, err := cache.GetClient()
	if err != nil {
		return false, err
	}
	cmd := client.SMove(cache.Name+src, cache.Name+dest, member)
	return cmd.Result()
}

// SPop 移除并返回集合的一个随机元素
func (cache *Cache) SPop(key string) (string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return "", err
	}
	cmd := client.SPop(cache.Name + key)
	return cmd.Result()
}

// SScan 迭代集合
func (cache *Cache) SScan(key string, cursor uint64, match string, count int64) ([]string, uint64) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, 0
	}
	cmd := client.SScan(cache.Name+key, cursor, match, count)
	return cmd.Val()
}

// HIncrBy 哈希自增
func (cache *Cache) HIncrBy(key, field string, incr int64) (int64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.HIncrBy(cache.Name+key, field, incr)
	return cmd.Result()
}

// HIncrByFloat 哈希自增
func (cache *Cache) HIncrByFloat(key, field string, incr float64) (float64, error) {
	client, err := cache.GetClient()
	if err != nil {
		return 0, err
	}
	cmd := client.HIncrByFloat(cache.Name+key, field, incr)
	return cmd.Result()
}

// HSetNX 设置哈希值
func (cache *Cache) HSetNX(key, field string, value interface{}) (bool, error) {
	client, err := cache.GetClient()
	if err != nil {
		return false, err
	}
	cmd := client.HSetNX(cache.Name+key, field, value)
	return cmd.Result()
}

// HExists 判断哈希值是否存在
func (cache *Cache) HExists(key, field string) (bool, error) {
	client, err := cache.GetClient()
	if err != nil {
		return false, err
	}
	cmd := client.HExists(cache.Name+key, field)
	return cmd.Result()
}

// HVals 获取哈希的所有值
func (cache *Cache) HVals(key string) ([]string, error) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, err
	}
	cmd := client.HVals(cache.Name + key)
	return cmd.Val(), cmd.Err()
}

// HScan 迭代哈希
func (cache *Cache) HScan(key string, cursor uint64, match string, count int64) ([]string, uint64) {
	client, err := cache.GetClient()
	if err != nil {
		return nil, 0
	}
	cmd := client.HScan(cache.Name+key, cursor, match, count)
	return cmd.Val()
}
