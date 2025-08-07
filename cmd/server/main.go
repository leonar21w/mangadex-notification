package main

import (
	"context"
	"log"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/constants"
	"github.com/leonar21w/mangadex-server-backend/internal/db"
	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/repository"
	"github.com/leonar21w/mangadex-server-backend/internal/services"
	"github.com/leonar21w/mangadex-server-backend/pkg"
)

func main() {
	cfg, err := pkg.Load()
	if err != nil {
		log.Fatalf("%v", err)
	}

	rdb, err := db.RedisInit(cfg.RedisURL, cfg.RedisToken)
	if err != nil {
		log.Fatalf("redis failed to initialize : %v", err)
	}

	//repositories
	tokenRepo := repository.NewRedisDB(rdb)
	mangaRepo := repository.NewRedisDB(rdb)

	//services
	authService := services.NewAuthService(tokenRepo)
	mangadexService := services.NewMangadexService(authService, mangaRepo)

	ctx := context.Background()
	client := models.UserAuthCredentials{
		GrantTye:     constants.MD.UserAuthGrantType(),
		Username:     cfg.Username,
		Password:     cfg.Password,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
	}

	tokenRepo.UpdateLastGetFeedTime(ctx)

	if err := authService.LoginWithMDX(ctx, client); err != nil {
		panic(err)
	}

	if err := mangadexService.InitializeMangas(ctx); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(constants.LoginInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := authService.LoginWithMDX(ctx, client); err != nil {
				log.Printf("relogin to mangadex error: %v", err)
			}
			log.Printf("logged in with %v", client.Username)
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.AccessTokenRefreshInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := authService.RefreshAccessTokens(ctx); err != nil {
				log.Printf("refresh tokens error: %v", err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.FollowedMangaListenerInterval)
		defer ticker.Stop()

		for {
			<-ticker.C
			//if no added mangas it nils out
			_, _, err := mangadexService.FetchMangasForAllClients(ctx)
			if err != nil {
				log.Print(err)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(constants.FeedListenerInterval)
		defer ticker.Stop()

		for {
			if err := mangadexService.AllClientsChapterFeed(ctx); err != nil {
				log.Printf("fetch notifications error: %v", err)
			}
			<-ticker.C
		}
	}()

	select {}
}
