package infra

import (
	"context"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClientOnce sync.Once
	redisClient     *redis.Client
)

func LoadRedis() *redis.Client {
	redisClientOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     LoadConfig().Redis.Address,
			Username: LoadConfig().Redis.Username,
			Password: LoadConfig().Redis.Password,
			DB:       LoadConfig().Redis.DB,
		})

		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("failed to ping redis: %v", err)
		}
	})

	return redisClient
}
