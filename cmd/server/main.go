package main

import (
	"log"
	"net/http"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/api"
	"github.com/leonar21w/mangadex-server-backend/internal/db"
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

	//handlers
	allHandlers := &api.Handler{
		Auth:     authService,
		Mangadex: mangadexService,
	}

	handler := api.NewHandler(allHandlers)

	router := api.NewChiRouter(handler)
	server := &http.Server{
		Addr:         ":5173",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
