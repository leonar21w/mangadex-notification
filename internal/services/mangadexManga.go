package services

import (
	"context"
	"net/http"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/util"
)

type MangadexService struct {
	TokensRepo   models.TokensRepo
	MangadexRepo models.MangadexRepo
}

func NewMangadexService(tokensRepo models.TokensRepo, mangadexRepo models.MangadexRepo) *MangadexService {
	return &MangadexService{
		TokensRepo:   tokensRepo,
		MangadexRepo: mangadexRepo,
	}
}

func (ms *MangadexService) FindAllMangasFollowedBy(ctx context.Context, clientID string) (*models.FollowedMangaCollection, error) {
	accessToken, err := ms.TokensRepo.GetAccessToken(ctx, clientID)
	if err != nil {
		return nil, err
	}

	endpoint := constants.MangaDexAPIBaseURL + "/user/follows/manga"
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	request, err := util.MakeHTTPRequest(ctx, endpoint, http.MethodGet, headers, nil, nil, models.FollowedMangaCollection{})
	if err != nil {
		return nil, err
	}
	return &request, nil
}
