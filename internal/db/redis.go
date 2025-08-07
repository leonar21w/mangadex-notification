package db

import (
	"context"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/redis/go-redis/v9"
)

func RedisInit(url string, password string) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	opt.Password = password
	opt.PoolSize = constants.Redis.PoolSize()
	opt.DialTimeout = constants.Redis.DialTimeout()
	opt.ReadTimeout = constants.Redis.ReadTimeout()
	opt.WriteTimeout = constants.Redis.WriteTimeout()

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), constants.Redis.ReadTimeout())

	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
