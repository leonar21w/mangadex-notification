package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leonar21w/mangadex-server-backend/internal/models"
	"github.com/leonar21w/mangadex-server-backend/internal/services"
)

// globals
type Handler struct {
	Auth *services.AuthService
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
