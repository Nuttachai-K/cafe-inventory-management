package handler

import (
	"errors"
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

func writeError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, service.ErrDataNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
