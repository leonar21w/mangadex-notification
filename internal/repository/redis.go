package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func (r *RedisDB) GetAllClients(ctx context.Context) (*models.ClientCollection, error) {
	rawClients, err := r.rdb.SMembers(ctx, "clients:mangadex").Result()
	if err != nil {
		return nil, err
	}

	var clients []models.Client
	for _, raw := range rawClients {
		parts := strings.SplitN(raw, ":", 2)
		if len(parts) != 2 {
			log.Printf("invalid client format in Redis: %s", raw)
			continue
		}
		client := models.Client{
			ClientID:     parts[0],
			ClientSecret: parts[1],
		}
		clients = append(clients, client)
	}

	return &models.ClientCollection{
		Clients: clients,
	}, nil
}

func (r *RedisDB) GetRefreshToken(ctx context.Context, clientID string) (string, error) {
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

func (r *RedisDB) CacheClientToken(ctx context.Context, t *models.Tokens, client *models.Client) error {
	rdb := r.rdb
	buildKeyAccess := fmt.Sprintf("access:%s", client.ClientID)
	buildKeyRefresh := fmt.Sprintf("refresh:%s", client.ClientID)

	if err := rdb.Set(ctx, buildKeyAccess, t.AccessToken, 10*time.Minute).Err(); err != nil {
		return err
	}
	if err := rdb.Set(ctx, buildKeyRefresh, t.RefreshToken, 24*28*time.Hour).Err(); err != nil {
		return err
	}
	buildClient := fmt.Sprintf("%s:%s", client.ClientID, client.ClientSecret)

	if err := rdb.SAdd(ctx, "clients:mangadex", buildClient).Err(); err != nil {
		return err
	}
	if err := rdb.Expire(ctx, "clients:mangadex", 24*28*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func (r *RedisDB) GetAccessToken(ctx context.Context, clientID string) (string, error) {
	rdb := r.rdb

	key := fmt.Sprintf("access:%s", clientID)

	accessToken, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}

	return accessToken, nil
}

func (r *RedisDB) CacheMangaIDList(ctx context.Context, mangaList []models.MangadexMangaData) error {
	for _, manga := range mangaList {
		if err := r.rdb.SAdd(ctx, "mangadex:mangaID", manga.ID).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (r *RedisDB) GetMangaIDList(ctx context.Context) ([]string, error) {
	mangaIDs, err := r.rdb.SMembers(ctx, "mangadex:mangaID").Result()
	if err != nil {
		return nil, err
	}
	return mangaIDs, nil
}

// should move to another interface
func (r *RedisDB) InsertMangaWithID(ctx context.Context, mangaID string, manga *models.Manga) error {
	key := fmt.Sprintf("mangadex:manga:%s", mangaID)

	if err := r.InsertAllChapters(ctx, mangaID, manga); err != nil {
		return err
	}
	return r.rdb.HSet(ctx, key, map[string]any{
		"title": manga.CanonicalTitle,
		"cover": manga.CoverURL,
	}).Err()
}

func (r *RedisDB) UpdateMangaField(ctx context.Context, mangaID string, field string, value any) error {
	key := fmt.Sprintf("mangadex:manga:%s", mangaID)
	return r.rdb.HSet(ctx, key, field, value).Err()
}

func (r *RedisDB) InsertAllChapters(ctx context.Context, mangaID string, manga *models.Manga) error {
	key := fmt.Sprintf("mangadex:manga:%s:chapters", mangaID)

	pipe := r.rdb.Pipeline()

	for i, chapter := range manga.Chapters {
		raw, err := json.Marshal(chapter)
		if err != nil {
			log.Print(i)
			return err
		}
		pipe.HSet(ctx, key, chapter.ID, raw)
	}
	_, err := pipe.Exec(ctx)
	return err
}
