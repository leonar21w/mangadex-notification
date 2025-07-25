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

func (r *RedisDB) GetAllCLients(ctx context.Context) ([]string, error) {
	rdb := r.rdb

	allMangadexClients, err := rdb.SMembers(ctx, "clients:mangadex").Result()
	if err != nil {
		return nil, err
	}

	return allMangadexClients, nil
}

func (r *RedisDB) GetRefreshTokens(ctx context.Context, clientID string) (string, error) {
	rdb := r.rdb

	buildKeyRefresh := fmt.Sprintf("refresh:%s", clientID)

	refreshToken, err := rdb.Get(ctx, buildKeyRefresh).Result()
	if err != nil {
		return "", fmt.Errorf("error getting refresh tokens, %v", err)
	}

	return refreshToken, nil
}

// used when access token reaches ttl
func (r *RedisDB) CacheAccessToken(ctx context.Context, accessToken string, clientID string) error {
	rdb := r.rdb

	buildKeyAccess := fmt.Sprintf("access:%s", clientID)

	if err := rdb.Set(ctx, buildKeyAccess, accessToken, 10*time.Minute).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisDB) CacheTokens(ctx context.Context, t *models.Tokens, clientID string) error {
	rdb := r.rdb
	buildKeyAccess := fmt.Sprintf("access:%s", clientID)
	buildKeyRefresh := fmt.Sprintf("refresh:%s", clientID)

	if err := rdb.Set(ctx, buildKeyAccess, t.AccessToken, 10*time.Minute).Err(); err != nil {
		return err
	}
	if err := rdb.Set(ctx, buildKeyRefresh, t.RefreshToken, 24*28*time.Hour).Err(); err != nil {
		return err
	}
	if err := rdb.SAdd(ctx, "clients:mangadex", clientID).Err(); err != nil {
		return err
	}
	if err := rdb.Expire(ctx, "clients:mangadex", 24*28*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) GetAllAvailableMangadexTokens(ctx context.Context, tokenKeyType string) ([]string, error) {
	rdb := r.rdb
	ids, err := rdb.SMembers(ctx, "clients:mangadex").Result()
	if err != nil {
		return nil, err
	}

	keys := make([]string, len(ids))
	for i, id := range ids {
		keys[i] = tokenKeyType + id
	}

	v, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	accessTokens := make([]string, 0, len(v))
	for _, token := range v {
		if token == redis.Nil {
			//if data here is ever nil then the database holds bad integrity, just throw an error
			//we detect the error and we can just check for refresh. If THAT errors then its a real error.
			return nil, nil
		}
		if result, ok := token.(string); ok {
			accessTokens = append(accessTokens, result)
		}
	}
	return accessTokens, nil
}
