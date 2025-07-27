package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/util"
)

type MangadexService struct {
	Auth *AuthService
}

func NewMangadexService(authService *AuthService) *MangadexService {
	return &MangadexService{
		Auth: authService,
	}
}

// Fetch this in a long interval (10 minutes - 30 minutes)
func (ms *MangadexService) FetchMangasForAllClients(ctx context.Context) ([]models.FollowedMangaCollection, error) {
	clients, err := ms.Auth.tokenRepo.GetAllClients(ctx)
	if err != nil {
		return nil, err
	}

	allClientsMangaCollection := make([]models.FollowedMangaCollection, len(clients.Clients))
	var wg sync.WaitGroup
	var errors []error

	for _, client := range clients.Clients {
		wg.Add(1)
		go func(client models.Client) {
			defer wg.Done()
			mangas, err := ms.FindAllMangasFollowedBy(ctx, client.ClientID)
			if err != nil {
				errors = append(errors, err)
			}

			clientMangaCollection := models.FollowedMangaCollection{
				ClientID:        client.ClientID,
				MangaCollection: mangas,
			}
			if err := ms.Auth.tokenRepo.CacheMangaIDList(ctx, mangas); err != nil {
				errors = append(errors, err)
			}
			allClientsMangaCollection = append(allClientsMangaCollection, clientMangaCollection)
		}(client)
	}
	wg.Wait()
	if len(errors) != 0 {
		for _, foundError := range errors {
			log.Println(foundError)
		}
		return nil, fmt.Errorf("found %v errors in FetchMangasForAllClients()", len(errors))
	}

	return allClientsMangaCollection, nil
}

func (ms *MangadexService) FindAllMangasFollowedBy(ctx context.Context, clientID string) ([]models.MangadexMangaData, error) {
	accessToken, err := ms.Auth.tokenRepo.GetAccessToken(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if accessToken == "" {
		ms.Auth.RefreshAccessTokens(ctx)
		accessToken, err = ms.Auth.tokenRepo.GetAccessToken(ctx, clientID)
		if err != nil {
			return nil, err
		}
	}

	endpoint := constants.MangaDexAPIBaseURL + "/user/follows/manga"
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	request, err := util.MakeHTTPRequest(ctx, endpoint, http.MethodGet, headers, nil, nil, models.ClientFollowedMangaCollectionResponse{})
	if err != nil {
		return nil, err
	}
	return request.Data, nil
}
