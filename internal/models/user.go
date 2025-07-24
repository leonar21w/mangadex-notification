package models

import (
	"context"
	"time"
)

type User struct {
	ID       string
	Username string
}

type UserRepo interface {
	CreateUser(ctx context.Context, u *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
}

// Given by the Client in /login Handler. This will be in the POST body
type UserAuthCredentials struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// Returned by mangadex, store this in redis and db
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    time.Time
}
