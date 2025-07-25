package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
)

type AuthService struct {
	tokenRepo models.TokensRepo
}

func NewAuthService(tr models.TokensRepo) *AuthService {
	return &AuthService{tokenRepo: tr}
}

func (as *AuthService) LoginWithMDX(
	ctx context.Context,
	creds models.UserAuthCredentials) error {

	refreshToken, _ := as.tokenRepo.GetRefreshTokens(ctx, creds.ClientID)
	if refreshToken != "" {
		return nil
	}

	formVals := url.Values{
		"grant_type":    {"password"},
		"username":      {creds.Username},
		"password":      {creds.Password},
		"client_id":     {creds.ClientID},
		"client_secret": {creds.ClientSecret},
	}

	endpoint := "https://auth.mangadex.org/realms/mangadex/protocol/openid-connect/token"
	resp, err := http.PostForm(endpoint, formVals)
	if err != nil {
		return fmt.Errorf("error calling %s, %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s", data)
	}

	var responseTokens models.MangadexLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseTokens); err != nil {
		return fmt.Errorf("failed to decode in mdxlogin: %v", err)
	}

	tokens := models.Tokens{
		ClientID:     creds.ClientID,
		AccessToken:  responseTokens.AccessToken,
		RefreshToken: responseTokens.RefreshToken,
	}

	if err := as.tokenRepo.CacheTokens(ctx, &tokens, creds.ClientID); err != nil {
		return err
	}

	return nil
}
