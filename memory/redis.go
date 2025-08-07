package memory

import (
	"context"
	"errors"
	"time"

	config2 "github.com/manjada/com/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

type RedisWrap struct {
}

type RedisInterface interface {
	Set(ctx context.Context, key string, data interface{}, duration *time.Duration) error
	GetString(ctx context.Context, key string) string
	HashSet(ctx context.Context, key string, data map[string]interface{}, duration *time.Duration) error
	HashGet(ctx context.Context, key string) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	GetInt(ctx context.Context, key string) (int, error)
	GetBoolean(ctx context.Context, key string) (bool, error)
	Increment(ctx context.Context, key string, duration ...time.Duration) (int64, error)
}

func (r RedisWrap) Increment(ctx context.Context, key string, duration ...time.Duration) (int64, error) {
	val, err := redisClient.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration if duration is provided
	if len(duration) > 0 && duration[0] > 0 {
		err = redisClient.PExpire(ctx, key, duration[0]).Err()
		if err != nil {
			return 0, err
		}
	}

	return val, nil
}

func NewRedisWrap() (*RedisWrap, error) {
	if redisClient == nil {
		config := config2.GetConfig().Redis
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.Address,
			Password: config.Pass,  // no password set
			DB:       config.Index, // use default DB
		})
	}
	if redisClient == nil {
		return nil, errors.New("Redis can't connected ")
	}
	return &RedisWrap{}, nil
}

func (r RedisWrap) Delete(ctx context.Context, key string) error {
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisWrap) Set(ctx context.Context, key string, data interface{}, duration *time.Duration) error {
	if duration == nil {
		*duration = 0
	}

	err := redisClient.Set(ctx, key, data, *duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisWrap) GetString(ctx context.Context, key string) string {
	val := redisClient.Get(ctx, key).Val()
	return val
}

func (r RedisWrap) GetBoolean(ctx context.Context, key string) (bool, error) {
	val, err := redisClient.Get(ctx, key).Bool()
	if err != nil {
		return false, err
	}
	return val, nil
}

func (r RedisWrap) GetInt(ctx context.Context, key string) (int, error) {
	val, err := redisClient.Get(ctx, key).Int()
	if err != nil {
		return 0, err
	}
	return val, nil
}

func (r RedisWrap) HashSet(ctx context.Context, key string, data map[string]interface{}, duration *time.Duration) error {
	for k, v := range data {
		err := redisClient.HSet(ctx, key, k, v).Err()
		if err != nil {
			config2.Error(err)
			return err
		}
	}
	return nil
}

func (r RedisWrap) HashGet(ctx context.Context, key string) (map[string]string, error) {
	userSession := redisClient.HGetAll(ctx, key).Val()
	return userSession, nil
}
