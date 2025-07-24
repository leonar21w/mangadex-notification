package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	rdb *redis.Client
}

func NewRedisDB(rdb *redis.Client) *RedisDB {
	return &RedisDB{rdb: rdb}
}

func (r *RedisDB) CacheTokens(ctx context.Context, t *models.Tokens, clientID string) error {
	rdb := r.rdb
	buildKeyAccess := fmt.Sprintf("access:%s", clientID)
	buildKeyRefresh := fmt.Sprintf("refresh:%s", clientID)

	rdb.Set(ctx, buildKeyAccess, t.AccessToken, 10*time.Minute)
	rdb.Set(ctx, buildKeyRefresh, t.RefreshToken, 24*28*time.Hour)
	rdb.SAdd(ctx, "clients:mangadex", clientID)
	return nil
}

func (r *RedisDB) GetAllAvailableMangadexAccessTokens(ctx context.Context) ([]string, error) {
	rdb := r.rdb
	ids, err := rdb.SMembers(ctx, "clients:mangadex").Result()
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(ids))

	for i, id := range ids {
		keys[i] = "access:" + id
	}

	v, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	accessTokens := make([]string, 0, len(v))
	for _, token := range v {
		if token == redis.Nil {
			//do refresh token
		}
		if result, ok := token.(string); ok {
			accessTokens = append(accessTokens, result)
		}
	}
	return accessTokens, nil
}
