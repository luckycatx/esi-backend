package db

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, opts *redis.Options) *redis.Client {
	var redis_client = redis.NewClient(opts)
	_, err := redis_client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return redis_client
}
