package cache

import (
	"auth-microservice/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, logger.Error(err, "failed to connect to redis")
	}

	return &RedisCache{
		client: client,
	}, nil
}

func (r *RedisCache) SetExpire(key string, value string, expiration time.Duration) error {
	err := r.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		return logger.Error(err, fmt.Sprintf("failed to set key %s", key))
	}
	return nil
}

func (r *RedisCache) Get(key string) (string, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if errors.Is(err, redis.Nil) {
		return "", logger.Error(err, fmt.Sprintf("key %s does not exist", key))
	} else if err != nil {
		return "", logger.Error(err, fmt.Sprintf("failed to get key %s", key))
	}
	return val, nil
}

func (r *RedisCache) Delete(key string) error {
	err := r.client.Del(context.Background(), key).Err()
	if err != nil {
		return logger.Error(err, fmt.Sprintf("failed to delete key %s", key))
	}
	return nil
}
