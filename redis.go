package mjd

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

type RedisWrap struct {
}

type RedisInterface interface {
	NewRedisWrap() (*RedisWrap, error)
	Set(ctx context.Context, key string, data interface{}, duration *time.Duration) error
	Get(ctx context.Context, key string) (interface{}, error)
	HashSet(ctx context.Context, key string, data map[string]interface{}, duration *time.Duration) error
	HashGet(ctx context.Context, key string) (map[string]string, error)
}

func (r RedisWrap) NewRedisWrap() (*RedisWrap, error) {
	if redisClient == nil {
		config := GetConfig().Redis
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

func (r RedisWrap) Set(ctx context.Context, key string, data interface{}, duration *time.Duration) error {
	if duration == nil {
		*duration = 0
	}
	err := redisClient.Set(ctx, key, data, *duration).Err()
	if err != nil {
		Error(err)
		return err
	}
	return nil
}

func (r RedisWrap) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (r RedisWrap) HashSet(ctx context.Context, key string, data map[string]interface{}, duration *time.Duration) error {
	for k, v := range data {
		err := redisClient.HSet(ctx, key, k, v).Err()
		if err != nil {
			Error(err)
			return err
		}
	}
	return nil
}

func (r RedisWrap) HashGet(ctx context.Context, key string) (map[string]string, error) {
	userSession := redisClient.HGetAll(ctx, key).Val()
	return userSession, nil
}
