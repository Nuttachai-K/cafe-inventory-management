package handler

import (
	"encoding/json"
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

// GwtCafeByID godoc
// @Summary Get a cafe by id
// @Tags cafes
// @Accept json
// @Produce json
// @Param id path int true "Cafe ID"
// @Success 200 {object} model.Cafe
// @Failure 400 {string} string "invalid cafe id"
// @Failure 404 {string} string "data not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/cafes/{id} [get]
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

// GetAllCafe godoc
// @Summary Get all cafes
// @Tags cafes
// @Accept json
// @Produce json
// @Param station query string false "Filter by nearest station"
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.Cafe
// @Failure 400 {string} string "invalid limit"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/cafes [get]
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

// CreateCafe godoc
// @Summary Create a cafe
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param cafe body model.Cafe true "Cafe payload"
// @Success 201 {object} object{id=int,message=string}
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/cafes [post]
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

// UpdateCafe godoc
// @Summary Update a cafe
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cafe ID"
// @Param cafe body model.CafeUpdate true "Cafe payload"
// @Success 200 {object} model.Cafe
// @Failure 400 {string} string "invalid request body or invalid cafe id"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/cafes/{id} [patch]
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

	cafe, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cafe); err != nil {
		writeError(w, err)
	}
}

// DeleteCafe godoc
// @Summary Delete a cafe
// @Tags cafes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Cafe ID"
// @Success 204
// @Failure 400 {string} string "invalid id"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "cafe has dependent records"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/cafes/{id} [delete]
func (h *CafeHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
