package common

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds connection settings for Redis.
type RedisConfig struct {
	Addr        string
	Password    string
	DB          int
	PingTimeout time.Duration
}

// OpenRedis creates a Redis client using the provided configuration.
func OpenRedis(cfg RedisConfig) (*redis.Client, error) {
	if cfg.PingTimeout == 0 {
		cfg.PingTimeout = 5 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}
