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
	ExpiresAt    time.Time
}

type Clients struct {
	ClientIDs []string
}

type TokensRepo interface {
	CacheTokens(ctx context.Context, t *Tokens, clientID string) error
	GetAllAvailableMangadexTokens(ctx context.Context, tokenKeyType string) ([]string, error)
	GetAllCLients(ctx context.Context) ([]string, error)
}
