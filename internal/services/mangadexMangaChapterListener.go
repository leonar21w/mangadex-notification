package services

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/util"
)

//This file is for scheduled tasks (1 mnt interval for chapters)

func (ms *MangadexService) AllClientsChapterFeed(ctx context.Context) error {
	clients, err := ms.Auth.tokenRepo.GetAllClients(ctx)
	if err != nil {
		return err
	}

	//MangaIDs : Chapter.Attributes
	mangaUpdates := make(map[string][]models.FeedChapter)
	var clientWg sync.WaitGroup
	var errors []error

	for _, client := range clients.Clients {
		clientWg.Add(1)
		go func(clientID string) {
			defer clientWg.Done()

			feed, err := ms.MangadexChapterFeed(ctx, clientID)
			if err != nil {
				errors = append(errors, err)
			}
			for _, chapterUpdates := range feed {
				if chapterUpdates.Attributes.TranslatedLanguage != "en" && chapterUpdates.Attributes.TranslatedLanguage != "id" {
					continue
				}
				lastTime, err := ms.Auth.tokenRepo.GetLastFeedTime(ctx)
				if err != nil {
					errors = append(errors, err)
				}
				parsedLastTime, err := time.ParseInLocation("2006-01-02T15:04:05", lastTime, time.UTC)
				if err != nil {
					errors = append(errors, err)
				}

				parsedChapterCreatedTime, err := time.Parse(time.RFC3339, chapterUpdates.Attributes.CreatedAt)
				if err != nil {
					errors = append(errors, err)
				}

				for _, rel := range chapterUpdates.Relationships {
					if rel.Type == "manga" && parsedChapterCreatedTime.After(parsedLastTime) {
						mangaUpdates[rel.ID] = append(mangaUpdates[rel.ID], chapterUpdates)
					}
				}
			}
		}(client.ClientID)
	}
	clientWg.Wait()

	var mangaWg sync.WaitGroup

	for mangaID, chapterUpdate := range mangaUpdates {
		mangaWg.Add(1)
		go func(mangaID string, chapterUpdate []models.FeedChapter) {
			defer mangaWg.Done()
			savedChapters, err := ms.Auth.tokenRepo.UpdateMangaChapters(ctx, mangaID, chapterUpdate)
			if err != nil {
				errors = append(errors, err)
			}

			if len(savedChapters) < 1 {
				return
			}

			log.Println(len(savedChapters), mangaID)

			title, _ := ms.Auth.tokenRepo.GetMangaTitle(ctx, mangaID)
			if title == "" {
				title = "no title found"
			}
			if err := SendMangaUpdateEmail(MangaUpdateEmailData{
				MangaTitle: title,
				MangaURL:   mangaID,
				Chapters:   savedChapters,
			}); err != nil {
				errors = append(errors, err)
			}
		}(mangaID, chapterUpdate)
	}
	mangaWg.Wait()

	if len(errors) > 0 {
		return errors[len(errors)-1]
	}

	return nil
}

func (ms *MangadexService) MangadexChapterFeed(ctx context.Context, clientID string) ([]models.FeedChapter, error) {
	accessToken, err := ms.Auth.tokenRepo.GetAccessToken(ctx, clientID)
	if err != nil {
		return nil, err
	}

	if accessToken == "" {
		if err := ms.Auth.RefreshAccessTokens(ctx); err != nil {
			return nil, err
		}

		accessToken, _ = ms.Auth.tokenRepo.GetAccessToken(ctx, clientID)
	}

	endpoint := constants.MangaDexAPIBaseURL + "/user/follows/manga/feed"
	header := map[string]string{
		"Authorization": "Bearer " + accessToken,
	}

	oldTimeMark, _ := ms.Auth.tokenRepo.GetLastFeedTime(ctx)

	if oldTimeMark == "" {
		// fallback to 1 minute ago
		oldTimeMark = time.Now().UTC().Add(-time.Minute).Format("2006-01-02T15:04:05")

	}
	limit := 100
	offset := 0

	fetchFeedPage := func(offset int) (models.FeedResponse, error) {
		queryParameters := url.Values{
			"limit":          {strconv.Itoa(limit)},
			"offset":         {strconv.Itoa(offset)},
			"publishAtSince": {oldTimeMark},
		}
		return util.MakeHTTPRequest(ctx, endpoint, string(http.MethodGet), header, queryParameters, nil, models.FeedResponse{})
	}

	firstPage, err := fetchFeedPage(offset)
	if err != nil {
		return nil, err
	}

	allChapters := make([]models.FeedChapter, 0, firstPage.Total)
	allChapters = append(allChapters, firstPage.Data...)

	pageCount := (((firstPage.Total) + limit - 1) / limit)

	for page := 1; page < pageCount; page++ {
		offset := page * limit
		nextPage, err := fetchFeedPage(offset)
		if err != nil {
			return nil, err
		}
		allChapters = append(allChapters, nextPage.Data...)
	}

	if err := ms.Auth.tokenRepo.UpdateLastGetFeedTime(ctx); err != nil {
		log.Println(err)
	}

	return allChapters, nil
}
