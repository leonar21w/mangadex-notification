package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
)

type AuthService struct {
	tokenRepo models.TokensRepo
}

func NewAuthService(tr models.TokensRepo) *AuthService {
	return &AuthService{tokenRepo: tr}
}

const authEndpoint = "/realms/mangadex/protocol/openid-connect/token"

func (as *AuthService) LoginWithMDX(
	ctx context.Context,
	creds models.UserAuthCredentials) error {

	refreshToken, _ := as.tokenRepo.GetRefreshToken(ctx, creds.ClientID)
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

	endpoint := constants.MangadexAuthBaseURL + authEndpoint
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

	client := models.Client{
		ClientID:     creds.ClientID,
		ClientSecret: creds.ClientSecret,
	}

	if err := as.tokenRepo.CacheClientToken(ctx, &tokens, &client); err != nil {
		return err
	}

	return nil
}

func (as *AuthService) RefreshAccessTokens(ctx context.Context) error {
	allClients, err := as.tokenRepo.GetAllClients(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errors []error

	for _, val := range allClients.Clients {
		wg.Add(1)
		go func(val models.Client) {
			defer wg.Done()

			refreshToken, err := as.tokenRepo.GetRefreshToken(ctx, val.ClientID)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			formVals := url.Values{
				"grant_type":    {"refresh_token"},
				"refresh_token": {refreshToken},
				"client_id":     {val.ClientID},
				"client_secret": {val.ClientSecret},
			}

			endpoint := constants.MangadexAuthBaseURL + authEndpoint
			resp, err := http.PostForm(endpoint, formVals)
			if err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				mu.Lock()
				data, _ := io.ReadAll(resp.Body)
				errors = append(errors, fmt.Errorf("unable to refresh this token: %s", val.ClientID))
				errors = append(errors, fmt.Errorf("received error while refreshing token: %v", data))
				mu.Unlock()
				return
			}

			var result models.RefreshTokenResponse

			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}

			if err := as.tokenRepo.CacheAccessToken(ctx, result.AccessToken, val.ClientID); err != nil {
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
		}(val)
	}
	wg.Wait()
	if len(errors) > 0 {
		return fmt.Errorf("encountered %d errors during token refresh, last error: %w", len(errors), errors[len(errors)-1])
	}

	return nil
}
