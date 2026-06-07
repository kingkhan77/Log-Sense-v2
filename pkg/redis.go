package pkg

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func NewRedis(cfg *Config) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			cfg.Redis.Host,
			cfg.Redis.Port,
		),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	return client
}