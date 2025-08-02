package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

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

	rdb, err := db.RedisInit(cfg.RedisURL, cfg.RedisToken, 5)
	if err != nil {
		log.Fatalf("redis failed to initialize : %v", err)
	}

	//repositories
	tokenRepo := repository.NewRedisDB(rdb)

	//services
	authService := services.NewAuthService(tokenRepo)
	mangadexService := services.NewMangadexService(authService)

	ctx := context.Background()
	client := models.UserAuthCredentials{
		GrantTye:     "password",
		Username:     os.Getenv("MGDEX_USERNAME"),
		Password:     os.Getenv("MGDEX_PASSWORD"),
		ClientID:     os.Getenv("MGDEX_CLIENT"),
		ClientSecret: os.Getenv("MGDEX_SECRET"),
	}

	if err := authService.LoginWithMDX(ctx, client); err != nil {
		panic(err)
	}

	if err := mangadexService.InitializeMangas(ctx); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(28 * 24 * time.Hour)
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
		ticker := time.NewTicker(20 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := authService.RefreshAccessTokens(ctx); err != nil {
				log.Printf("refresh tokens error: %v", err)
			}
			log.Printf("refreshed tokens")
		}
	}()

	go func() {
		ticker := time.NewTicker(3 * time.Hour)
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
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			if err := mangadexService.AllClientsChapterFeed(ctx); err != nil {
				log.Printf("fetch notifications error: %v", err)
			}
			log.Printf("fetched 2 mnt interval")
			<-ticker.C
		}
	}()

	//Satisfying render server

	// Start the HTTP listener on Render's port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	go func() {
		log.Printf("🌐 Listening for health checks on :%s", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	select {}
}
