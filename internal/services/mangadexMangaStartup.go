package services

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/util"
)

//This file contains behaviour for fetching mangas on database startup.

func (ms *MangadexService) FetchMangasChapters(ctx context.Context, mangaData *models.MangadexMangaData) (*models.Manga, error) {
	chapterList, err := ms.FetchAllChapterList(ctx, mangaData.ID)
	if err != nil {
		return nil, err
	}
	title, ok := mangaData.Attributes.Title["en"]
	if !ok {
		title, ok = mangaData.Attributes.Title["jp"]
	}
	if !ok {
		title, ok = mangaData.Attributes.Title["id"]
	}
	if !ok {
		title = "title not available"
	}

	return &models.Manga{
		ID:             mangaData.ID,
		CanonicalTitle: title,
		Chapters:       chapterList.Chapters,
		CoverURL:       mangaData.Attributes.CoverURL,
	}, nil
}

func (ms *MangadexService) FetchAllChapterList(ctx context.Context, mangaID string) (*models.MangadexChapterList, error) {
	endpoint := constants.MangaDexAPIBaseURL + "/chapter"

	header := map[string]string{
		"Content-Type": "application/json",
	}

	var allChapters []models.MangadexChapterData
	offset := 0
	limit := 100

	for {
		queryParam := url.Values{
			"manga":                {mangaID},
			"translatedLanguage[]": {"en"},
			"order[chapter]":       {"asc"},
			"limit":                {fmt.Sprintf("%d", limit)},
			"offset":               {fmt.Sprintf("%d", offset)},
		}

		res, err := util.MakeHTTPRequest(ctx, endpoint, http.MethodGet, header, queryParam, nil, models.MangadexChapterListResponse{})
		if err != nil {
			return nil, err
		}

		// Accumulate chapters
		allChapters = append(allChapters, res.Data...)

		// Break if all data is fetched
		if offset+limit >= res.Total {
			break
		}

		offset += limit
	}

	return &models.MangadexChapterList{
		Chapters: allChapters,
	}, nil
}
