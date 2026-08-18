package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

func writeError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeJSONError(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, service.ErrDataNotFound):
		writeJSONError(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, service.ErrDuplicateEmail):
		writeJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrDuplicateCategory):
		writeJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrInvalidCredentials):
		writeJSONError(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, service.ErrInsufficientStock):
		writeJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrHasDependents):
		writeJSONError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, service.ErrUserInactive):
		writeJSONError(w, err.Error(), http.StatusUnauthorized)

	default:
		writeJSONError(w, "internal server error", http.StatusInternalServerError)
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
