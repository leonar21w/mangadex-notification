package models

import (
	"context"
	"time"
)

// Returned by mangadex, store this in redis and db
type UserAuthCredentials struct {
	GrantTye     string
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type Tokens struct {
	ClientID     string
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MangadexLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    time.Time
}

type ClientCollection struct {
	Clients []Client
}

type Client struct {
	ClientID     string
	ClientSecret string
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type TokensRepo interface {
	CacheMangaIDList(ctx context.Context, mangaID []MangadexMangaData) error
	GetAllClients(ctx context.Context) (*ClientCollection, error)
	GetRefreshToken(ctx context.Context, clientID string) (string, error)
	GetAccessToken(ctx context.Context, clienID string) (string, error)
	CacheAccessToken(ctx context.Context, accessToken string, clientID string) error
	CacheClientToken(ctx context.Context, t *Tokens, client *Client) error
}
