package mjd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

var redisClient *redis.Client

type RedisWrap struct {
}

type RedisInterface interface {
	Set(ctx context.Context, key string, data interface{}, duration *time.Duration) error
	GetString(ctx context.Context, key string) (string, error)
	HashSet(ctx context.Context, key string, data map[string]interface{}, duration *time.Duration) error
	HashGet(ctx context.Context, key string) (map[string]string, error)
	Delete(ctx context.Context, key string) error
	GetInt(ctx context.Context, key string) (int, error)
}

func NewRedisWrap() (*RedisWrap, error) {
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
	dataByte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = redisClient.Set(ctx, key, string(dataByte), *duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r RedisWrap) GetString(ctx context.Context, key string) (string, error) {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
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
