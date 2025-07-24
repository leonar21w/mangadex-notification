package main

import (
	"log"
	"net/http"
	"time"

	"github.com/leonar21w/mangadex-server-backend/internal/api"
)

func main() {

	router := api.NewChiRouter()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Minute,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
