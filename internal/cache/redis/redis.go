package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"goods-manager/internal/cache"
	"time"
)

// Cache implementation `cache.Cache` using redis
type Cache struct {
	client *redis.Client
}

func (c Cache) Get(ctx context.Context, key string, value any) error {
	res, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return cache.ErrorNotExists
		}

		return err
	}

	return json.Unmarshal(res, value)
}

// Set data to store.
//
// Default ttl is 1 minute
func (c Cache) Set(ctx context.Context, key string, value any) error {
	valueJson, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, valueJson, 1*time.Minute).Err()
}

func (c Cache) Remove(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func NewCache(client *redis.Client) cache.Cache {
	return &Cache{client: client}
}
