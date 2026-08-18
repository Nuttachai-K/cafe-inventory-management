package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// GetCategoryByID godoc
// @Summary Get a category by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.Category
// @Failure 400 {string} string "invalid category id"
// @Failure 404 {string} string "data not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(
			w,
			"Invalid user id",
			http.StatusBadRequest,
		)
		return
	}

	category, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(category); err != nil {
		writeError(w, err)
	}
}

// GetAllCategory godoc
// @Summary Get all categories
// @Tags categories
// @Accept json
// @Produce json
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.Category
// @Failure 400 {string} string "invalid limit"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			writeJSONError(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	categories, err := h.service.GetAll(
		r.Context(),
		limit,
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(categories); err != nil {
		writeError(w, err)
	}
}

// CreateCategory godoc
// @Summary Create a category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body object{name=string} true "Category payload"
// @Success 201 {object} object{id=int,message=string}
// @Failure 400 {string} string "invalid request body"
// @Failure 401 {string} string "invalid or expired token"
// @Failure 403 {string} string "permission denied"
// @Failure 409 {string} string "category name is already used"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var category model.Category

	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &category); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      category.ID,
		Message: "Category created successfully",
	}); err != nil {
		writeError(w, err)
	}
}

// UpdateCategory godoc
// @Summary Update a category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param category body object{name=string} true "Category payload"
// @Success 200 {object} model.Category
// @Failure 400 {string} string "invalid request body or invalid category id"
// @Failure 401 {string} string "invalid or expired token"
// @Failure 403 {string} string "permission denied"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "category name is already used"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/categories/{id} [patch]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	category, err := h.service.Update(r.Context(), id, req.Name)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(category); err != nil {
		writeError(w, err)
	}
}

// DeleteCategory godoc
// @Summary Delete a category
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {string} string "invalid category id"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "category has dependent records"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
