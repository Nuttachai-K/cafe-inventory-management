package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type CafeHandler struct {
	service service.CafeService
}

func NewCafeHandler(service service.CafeService) *CafeHandler {
	return &CafeHandler{
		service: service,
	}
}

func (h *CafeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(
			w,
			"Invalid cafe id",
			http.StatusBadRequest,
		)
		return
	}

	cafe, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cafe); err != nil {
		writeError(w, err)
	}
}

func (h *CafeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	station := r.URL.Query().Get("station")

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	cafes, err := h.service.GetAll(
		r.Context(),
		station,
		limit,
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cafes); err != nil {
		writeError(w, err)
	}
}

func (h *CafeHandler) Create(w http.ResponseWriter, r *http.Request) {

	var cafe model.Cafe

	if err := json.NewDecoder(r.Body).Decode(&cafe); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &cafe); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      cafe.ID,
		Message: "Cafe created successfully",
	}); err != nil {
		writeError(w, err)
	}

}

func (h *CafeHandler) Update(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.CafeUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), id, &req); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CafeHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, service.ErrCafeNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)

	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
