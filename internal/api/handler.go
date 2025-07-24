package api

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) PingPong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ponged you")
}

func (h *Handler) MDXLogin(w http.ResponseWriter, r *http.Request) {
	//Login, cache tokens in Redis, store UUID in database where UUID mapped to username

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}
}
