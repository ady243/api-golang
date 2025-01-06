package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(addr, password string, db int) *RedisService {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisService{
		client: client,
	}
}

func (rs *RedisService) SetNotification(key string, value string, expiration time.Duration) error {
	ctx := context.Background()
	return rs.client.Set(ctx, key, value, expiration).Err()
}

func (rs *RedisService) GetNotification(key string) (string, error) {
	ctx := context.Background()
	return rs.client.Get(ctx, key).Result()
}

func (rs *RedisService) DeleteNotification(key string) error {
	ctx := context.Background()
	return rs.client.Del(ctx, key).Err()
}

func (rs *RedisService) GetKeys(pattern string) ([]string, error) {
	ctx := context.Background()
	return rs.client.Keys(ctx, pattern).Result()
}
