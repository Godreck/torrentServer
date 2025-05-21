// storage/redis.go
package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		}),
		ttl: ttl,
	}
}

// Сохранить результат в кэш
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		log.Printf("Cache marshal error: %v (key: %s)", err, key)
		return err
	}
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		log.Printf("Cache set error: %v (key: %s)", err, key)
		return err
	}
	log.Printf("Cached data for key: %s", key)
	return nil
}

// Получить результат из кэша
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) bool {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Key %s not found in cache", key)
		} else {
			log.Printf("Cache get error: %v", err)
		}
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		log.Printf("Cache unmarshal error: %v", err)
		return false
	}
	return true
}
