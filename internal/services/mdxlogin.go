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

	allMangadexClients, err := as.tokenRepo.GetAllCLients(ctx)
	if err != nil {
		return err
	}

	isRegisteredMember := func(allClients []string, client string) bool {
		for _, item := range allClients {
			if item == client {
				return true
			}
		}
		return false
	}(allMangadexClients, creds.ClientID)

	if isRegisteredMember {
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

	var tokens models.Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return fmt.Errorf("failed to decode in mdxlogin: %v", err)
	}

	if err := as.tokenRepo.CacheTokens(ctx, &tokens, creds.ClientID); err != nil {
		return err
	}

	return nil
}
