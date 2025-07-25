package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewChiRouter(allHandlers *Handler) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)
	r.Get("/ping", allHandlers.PingPong)

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", allHandlers.Login)
		r.Get("/refresh", allHandlers.RefreshMangadexAccessTokens)
	})

	return r
}
