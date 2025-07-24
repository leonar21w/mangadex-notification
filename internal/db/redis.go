package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func RedisInit(url string, password string, poolSize int) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	opt.Password = password
	opt.PoolSize = poolSize
	opt.DialTimeout = 6 * time.Second
	opt.ReadTimeout = 5 * time.Second
	opt.WriteTimeout = 5 * time.Second

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
