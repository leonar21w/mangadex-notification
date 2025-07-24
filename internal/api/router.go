package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewChiRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	h := NewHandler()
	r.Route("/v2/api", func(r chi.Router) {
		r.Post("/login", h.MDXLogin)
	})

	return r
}
