package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param login body model.Authentication true "Login payload"
// @Success 201 {object} object{token=string}
// @Failure 400 {object} model.ErrorResponse "invalid request body"
// @Failure 401 {object} model.ErrorResponse "invalid credential or inactive user"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var auth model.Authentication

	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), &auth)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Token string `json:"token"`
	}{
		Token: token,
	}); err != nil {
		writeError(w, err)
	}
}
