package storage

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// InitRedis создаёт Redis клиента и проверяет подключение
func InitRedis(ctx context.Context, opts *redis.Options) *redis.Client {
	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis connect error: %v", err)
	}
	return client
}