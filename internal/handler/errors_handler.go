package handler

import (
	"errors"
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

func writeError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, service.ErrDataNotFound):
		utils.WriteJSONError(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, service.ErrDuplicateEmail):
		utils.WriteJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrDuplicateCategory):
		utils.WriteJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrInvalidCredentials):
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, service.ErrInsufficientStock):
		utils.WriteJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrHasDependents):
		utils.WriteJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrUserInactive):
		utils.WriteJSONError(w, err.Error(), http.StatusUnauthorized)

	default:
		utils.WriteJSONError(w, "internal server error", http.StatusInternalServerError)
	}
}
