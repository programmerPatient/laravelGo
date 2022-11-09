package redis

import (
	"context"
	"sync"
	"time"

	"github.com/laravelGo/core/logger"

	redis "github.com/go-redis/redis/v8"
)

// RedisClient Redis 服务
type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// once 确保全局的 Redis 对象只实例一次
var once sync.Once

// Redis 全局 Redis，使用 db 1
var Redis *RedisClient

// ConnectRedis 连接 redis 数据库，设置全局的 Redis 对象
func ConnectRedis(address string, username string, password string, db int) {
	once.Do(func() {
		Redis = NewClient(address, username, password, db)
	})
}

// NewClient 创建一个新的 redis 连接
func NewClient(address string, username string, password string, db int) *RedisClient {
	// 初始化自定的 RedisClient 实例
	rds := &RedisClient{}
	// 使用默认的 context
	rds.Context = context.Background()

	// 使用 redis 库里的 NewClient 初始化连接
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       db,
	})

	// 测试一下连接
	err := rds.Ping()
	logger.LogIf(err)
	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
		logger.ErrorString("Redis", "Set", err.Error())
		return false
	}
	return true
}

// Get 获取 key 对应的 value
func (rds RedisClient) Get(key string) string {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Get", err.Error())
		}
		return ""
	}
	return result
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds RedisClient) Has(key string) bool {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Has", err.Error())
		}
		return false
	}
	return true
}

func (rds RedisClient) Exists(keys ...string) int64 {
	var result int64
	var err error
	if result, err = rds.Client.Exists(rds.Context, keys...).Result(); err != nil {
		logger.ErrorString("Redis", "Exists", err.Error())
	}
	return result
}

func (rds RedisClient) Expire(key string, expiration time.Duration) bool {
	var result bool
	var err error
	if result, err = rds.Client.Expire(rds.Context, key, expiration).Result(); err != nil {
		logger.ErrorString("Redis", "Expire", err.Error())
	}
	return result
}

func (rds RedisClient) ExpireNX(key string, expiration time.Duration) bool {
	var result bool
	var err error
	if result, err = rds.Client.ExpireNX(rds.Context, key, expiration).Result(); err != nil {
		logger.ErrorString("Redis", "ExpireNX", err.Error())
	}
	return result
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds RedisClient) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Context, keys...).Err(); err != nil {
		logger.ErrorString("Redis", "Del", err.Error())
		return false
	}
	return true
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds RedisClient) FlushDB() bool {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		logger.ErrorString("Redis", "FlushDB", err.Error())
		return false
	}
	return true
}

// Increment 当参数只有 1 个时，为 key，其值增加 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
func (rds RedisClient) Increment(parameters ...interface{}) (bool, int64) {
	var num int64
	var err error
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if num, err = rds.Client.Incr(rds.Context, key).Result(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false, 0
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if num, err = rds.Client.IncrBy(rds.Context, key, value).Result(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false, 0
		}
	default:
		logger.ErrorString("Redis", "Increment", "参数过多")
		return false, 0
	}
	return true, num
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
func (rds RedisClient) Decrement(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Decr(rds.Context, key).Err(); err != nil {
			logger.ErrorString("Redis", "Decrement", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.DecrBy(rds.Context, key, value).Err(); err != nil {
			logger.ErrorString("Redis", "Decrement", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Decrement", "参数过多")
		return false
	}
	return true
}

func (rds RedisClient) SAdd(key string, members ...interface{}) bool {
	if err := rds.Client.SAdd(rds.Context, key, members...).Err(); err != nil {
		logger.ErrorString("Redis", "SAdd", err.Error())
		return false
	}
	return true
}

func (rds RedisClient) SCard(key string) int64 {
	var count int64
	var err error
	if count, err = rds.Client.SCard(rds.Context, key).Result(); err != nil {
		logger.ErrorString("Redis", "SCard", err.Error())
		count = 0
	}
	return count
}

func (rds RedisClient) HDel(key string, fields ...string) int64 {
	var result int64
	var err error
	if result, err = rds.Client.HDel(rds.Context, key, fields...).Result(); err != nil {
		logger.ErrorString("Redis", "HDel", err.Error())
		result = 0
	}
	return result
}
func (rds RedisClient) HExists(key string, fields string) bool {
	var result bool
	var err error
	if result, err = rds.Client.HExists(rds.Context, key, fields).Result(); err != nil {
		logger.ErrorString("Redis", "HExists", err.Error())
	}
	return result
}
func (rds RedisClient) HGet(key, field string) string {
	var result string
	var err error
	if result, err = rds.Client.HGet(rds.Context, key, field).Result(); err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "HGet", err.Error())
		}
		result = ""
	}
	return result
}

func (rds RedisClient) HGetAll(key string) map[string]string {
	var result map[string]string
	var err error
	if result, err = rds.Client.HGetAll(rds.Context, key).Result(); err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "HGetAll", err.Error())
		}
		result = map[string]string{}
	}
	return result
}

func (rds RedisClient) HIncrBy(key string, field string, incr int64) int64 {
	var result int64
	var err error
	if result, err = rds.Client.HIncrBy(rds.Context, key, field, incr).Result(); err != nil {
		logger.ErrorString("Redis", "HIncrBy", err.Error())
	}
	return result
}

func (rds RedisClient) HIncrByFloat(key string, field string, incr float64) float64 {
	var result float64
	var err error
	if result, err = rds.Client.HIncrByFloat(rds.Context, key, field, incr).Result(); err != nil {
		logger.ErrorString("Redis", "HIncrByFloat", err.Error())
	}
	return result
}

func (rds RedisClient) HKeys(key string) []string {
	var result []string
	var err error
	if result, err = rds.Client.HKeys(rds.Context, key).Result(); err != nil {
		logger.ErrorString("Redis", "HKeys", err.Error())
	}
	return result
}

func (rds RedisClient) HLen(key string) int64 {
	var result int64
	var err error
	if result, err = rds.Client.HLen(rds.Context, key).Result(); err != nil {
		logger.ErrorString("Redis", "HLen", err.Error())
	}
	return result
}

func (rds RedisClient) HMSet(key string, values ...interface{}) bool {
	var result bool
	var err error
	if result, err = rds.Client.HMSet(rds.Context, key, values...).Result(); err != nil {
		logger.ErrorString("Redis", "HMSet", err.Error())
	}
	return result
}

func (rds RedisClient) HSet(key string, values ...interface{}) int64 {
	var result int64
	var err error
	if result, err = rds.Client.HSet(rds.Context, key, values...).Result(); err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "HSet", err.Error())
		}
		result = 0
	}
	return result
}

func (rds RedisClient) HSetNX(key, field string, value interface{}) bool {
	var result bool
	var err error
	if result, err = rds.Client.HSetNX(rds.Context, key, field, value).Result(); err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "HSetNX", err.Error())
		}
	}
	return result
}

func (rds RedisClient) HVals(key string) []string {
	var result []string
	var err error
	if result, err = rds.Client.HVals(rds.Context, key).Result(); err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "HVals", err.Error())
		}
	}
	return result
}

func (rds RedisClient) SIsMember(key string, member interface{}) bool {
	var result bool
	var err error
	if result, err = rds.Client.SIsMember(rds.Context, key, member).Result(); err != nil {
		logger.ErrorString("Redis", "SIsMember", err.Error())
	}
	return result
}

func (rds RedisClient) SMIsMember(key string, member ...interface{}) []bool {
	var result []bool
	var err error
	if result, err = rds.Client.SMIsMember(rds.Context, key, member...).Result(); err != nil {
		logger.ErrorString("Redis", "SMIsMember", err.Error())
	}
	return result
}
