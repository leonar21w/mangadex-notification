package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/services"
)

type Handler struct {
	Auth     *services.AuthService
	Mangadex *services.MangadexService
}

func NewHandler(h *Handler) *Handler {
	return h
}

func (h *Handler) PingPong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("ponged you")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		return
	}

	var credentials models.UserAuthCredentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err := h.Auth.LoginWithMDX(r.Context(), credentials); err != nil {
		errString := fmt.Errorf("invalid login with mangadex %v", err).Error()
		http.Error(w, errString, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successfull",
	})
}

func (h *Handler) RefreshMangadexAccessTokens(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		return
	}

	if err := h.Auth.RefreshAccessTokens(r.Context()); err != nil {
		json.NewEncoder(w).Encode("could not refresh access to mangadex")
		return
	}

	json.NewEncoder(w).Encode("success")
}

func (h *Handler) tryingendpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	someCollection, err := h.Mangadex.FetchMangasForAllClients(r.Context())

	if err != nil {
		json.NewEncoder(w).Encode(fmt.Sprintf("error found: %v", err))
		return
	}

	json.NewEncoder(w).Encode(someCollection)
}
