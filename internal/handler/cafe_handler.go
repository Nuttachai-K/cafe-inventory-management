package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

type CafeHandler struct {
	service service.CafeService
}

func NewCafeHandler(service service.CafeService) *CafeHandler {
	return &CafeHandler{
		service: service,
	}
}

// GetCafeByID godoc
// @Summary Get a cafe by id
// @Tags cafes
// @Accept json
// @Produce json
// @Param id path int true "Cafe ID"
// @Success 200 {object} model.Cafe
// @Failure 400 {object} model.ErrorResponse "invalid cafe id"
// @Failure 404 {object} model.ErrorResponse "data not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/cafes/{id} [get]
func (h *CafeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(
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
// @Param station query string false "Filter by nearest station (partial, case-insensitive match)"
// @Param lat query number false "Customer latitude, must be provided together with lng"
// @Param lng query number false "Customer longitude, must be provided together with lat"
// @Param radius query number false "Max distance in km from lat/lng (requires lat and lng)"
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.Cafe
// @Failure 400 {object} model.ErrorResponse "invalid query parameter"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/cafes [get]
func (h *CafeHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	var filter model.CafeFilter

	if v := r.URL.Query().Get("station"); v != "" {
		filter.Station = &v
	}

	if v := r.URL.Query().Get("lat"); v != "" {
		lat, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.WriteJSONError(w, "invalid latitude", http.StatusBadRequest)
			return
		}
		filter.Lat = &lat
	}

	if v := r.URL.Query().Get("lng"); v != "" {
		lng, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.WriteJSONError(w, "invalid longitude", http.StatusBadRequest)
			return
		}
		filter.Lng = &lng
	}

	if v := r.URL.Query().Get("radius"); v != "" {
		radius, err := strconv.ParseFloat(v, 64)
		if err != nil {
			utils.WriteJSONError(w, "invalid radius", http.StatusBadRequest)
			return
		}
		filter.RadiusKm = &radius
	}

	filter.Limit = 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			utils.WriteJSONError(w, "invalid limit", http.StatusBadRequest)
			return
		}
		filter.Limit = parsed
	}

	cafes, err := h.service.GetAll(
		r.Context(),
		filter,
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
// @Param cafe body model.CafeCreate true "Cafe payload"
// @Success 201 {object} object{id=int,message=string}
// @Failure 400 {object} model.ErrorResponse "invalid request body"
// @Failure 401 {object} model.ErrorResponse "invalid or expired token"
// @Failure 403 {object} model.ErrorResponse "permission denied"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/cafes [post]
func (h *CafeHandler) Create(w http.ResponseWriter, r *http.Request) {

	var cafe model.Cafe

	if err := json.NewDecoder(r.Body).Decode(&cafe); err != nil {
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
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
// @Failure 400 {object} model.ErrorResponse "invalid request body or invalid cafe id"
// @Failure 401 {object} model.ErrorResponse "invalid or expired token"
// @Failure 403 {object} model.ErrorResponse "permission denied"
// @Failure 404 {object} model.ErrorResponse "data not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/cafes/{id} [patch]
func (h *CafeHandler) Update(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.CafeUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
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
// @Failure 400 {object} model.ErrorResponse "invalid id"
// @Failure 401 {object} model.ErrorResponse "invalid or expired token"
// @Failure 403 {object} model.ErrorResponse "permission denied"
// @Failure 404 {object} model.ErrorResponse "data not found"
// @Failure 409 {object} model.ErrorResponse "cafe has dependent records"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/cafes/{id} [delete]
func (h *CafeHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
