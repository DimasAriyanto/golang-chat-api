package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(addr string, password string, db int) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *RedisCache) SetCache(key string, value string, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

func (c *RedisCache) GetCache(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

func (c *RedisCache) DeleteCache(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
