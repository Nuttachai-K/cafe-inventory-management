package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// GetProductByID godoc
// @Summary Get a product by id
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.ProductWithCategory
// @Failure 400 {string} string "invalid product id"
// @Failure 404 {string} string "data not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(
			w,
			"Invalid product id",
			http.StatusBadRequest,
		)
		return
	}

	product, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		writeError(w, err)
	}
}

// GetAllProduct godoc
// @Summary Get all products
// @Tags products
// @Accept json
// @Produce json
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.ProductWithCategory
// @Failure 400 {string} string "invalid limit"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	products, err := h.service.GetAll(
		r.Context(),
		limit,
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		writeError(w, err)
	}
}

// CreateProduct godoc
// @Summary Create a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product body model.ProductCreate true "Product payload"
// @Success 201 {object} object{id=int,message=string}
// @Failure 400 {string} string "invalid request body"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var product model.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &product); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      product.ID,
		Message: "Product created successfully",
	}); err != nil {
		writeError(w, err)
	}
}

// UpdateProduct godoc
// @Summary Update a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param product body model.ProductUpdate true "Product payload"
// @Success 200 {object} model.ProductWithCategory
// @Failure 400 {string} string "invalid request body or invalid product id"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "referenced cafe_id or category_id does not exist"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/products/{id} [patch]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.ProductUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		writeError(w, err)
	}
}

// DeleteProduct godoc
// @Summary Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {string} string "invalid product id"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "product has dependent records"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {

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
