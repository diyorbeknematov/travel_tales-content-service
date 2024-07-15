package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	R *redis.Client
}

var ctx = context.Background()

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &RedisClient{
		R: client,
	}
}

