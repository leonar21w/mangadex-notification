package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
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

// Needs manga clients to be in redis
func (ms *MangadexService) InitializeMangas(ctx context.Context) error {
	added, mangaList, err := ms.FetchMangasForAllClients(ctx)
	if err != nil {
		return err
	}
	if added == 0 {
		log.Print("no added mangas to cache")
		return nil
	}

	var wg sync.WaitGroup
	var errors []error

	for _, val := range mangaList {
		wg.Add(1)
		go func(val models.MangadexMangaData) {
			defer wg.Done()
			manga, err := ms.FetchMangasChapters(ctx, &val)
			if err != nil {
				errors = append(errors, err)
			}

			var wait sync.WaitGroup
			wait.Add(1)

			go func(valID string, manga *models.Manga) {
				defer wait.Done()
				ms.Auth.tokenRepo.InsertMangaWithID(ctx, val.ID, manga)

			}(val.ID, manga)

			ms.Auth.tokenRepo.InsertAllChapters(ctx, val.ID, manga)

			wait.Wait()
		}(val)
	}
	wg.Wait()

	if len(errors) > 0 {
		return errors[len(errors)-1]
	}
	return nil
}

// Fetch this in a long interval (10 minutes - 30 minutes)
func (ms *MangadexService) FetchMangasForAllClients(ctx context.Context) (int, []models.MangadexMangaData, error) {
	clients, err := ms.Auth.tokenRepo.GetAllClients(ctx)
	if err != nil {
		return 0, nil, err
	}

	var mangaList []models.MangadexMangaData
	var wg sync.WaitGroup
	var errors []error
	added := 0

	for _, client := range clients.Clients {
		wg.Add(1)
		go func(client models.Client) {
			defer wg.Done()
			mangas, err := ms.FindAllMangasFollowedBy(ctx, client.ClientID)
			if err != nil {
				errors = append(errors, err)
			}

			result, err := ms.Auth.tokenRepo.CacheMangaIDList(ctx, mangas)
			if err != nil {
				errors = append(errors, err)
			}
			added += result
			mangaList = append(mangaList, mangas...)
		}(client)
	}
	wg.Wait()
	if len(errors) != 0 {
		for _, foundError := range errors {
			log.Println(foundError)
		}
		return 0, nil, fmt.Errorf("found %v errors in FetchMangasForAllClients()", len(errors))
	}

	return added, mangaList, nil

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
	var all []models.MangadexMangaData
	endpoint := constants.MangaDexAPIBaseURL + "/user/follows/manga"
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}
	total := 1
	offset := 0

	for offset < total {
		params := url.Values{
			"limit":  {strconv.Itoa(constants.DefaultPageLimit)},
			"offset": {strconv.Itoa(offset)},
		}

		// fire the request
		req, err := util.MakeHTTPRequest(
			ctx,
			endpoint,
			http.MethodGet,
			headers,
			params,
			nil,
			models.ClientFollowedMangaCollectionResponse{},
		)
		if err != nil {
			return nil, err
		}

		// append this page
		all = append(all, req.Data...)

		total = req.Total
		offset += constants.DefaultPageLimit
	}

	return all, nil
}
