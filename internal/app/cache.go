package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// ConnectToRedis connect to redis
func ConnectToRedis(addr, pass string, db int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client, client.Ping(ctx).Err()
}
