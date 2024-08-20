package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	R *redis.Client
}

func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: "redis2:6379",
	})

	return &RedisClient{
		R: client,
	}
}
