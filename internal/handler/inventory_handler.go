package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/middleware"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		service: service,
	}
}

// GetInventoryByID godoc
// @Summary Get a inventory by id
// @Tags inventory
// @Accept json
// @Produce json
// @Param id path int true "Inventory ID"
// @Success 200 {object} model.InventoryWithProduct
// @Failure 400 {string} string "invalid inventory id"
// @Failure 404 {string} string "data not found"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/inventory/{id} [get]
func (h *InventoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {

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

	inventory, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		writeError(w, err)
	}
}

// GetAllInventory godoc
// @Summary Get all inventories
// @Tags inventory
// @Accept json
// @Produce json
// @Param name query string false "Product name"
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.InventoryWithProduct
// @Failure 400 {string} string "invalid limit"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/inventory [get]
func (h *InventoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	productName := r.URL.Query().Get("name")

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	inventories, err := h.service.GetAll(
		r.Context(),
		productName,
		limit,
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(inventories); err != nil {
		writeError(w, err)
	}
}

// UpdateStock godoc
// @Summary Update a stock
// @Tags inventory
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Product ID"
// @Param inventory body model.InventoryUpdateRequest true "Inventory payload"
// @Success 200 {object} object{message=string, stock_quantity=int}
// @Failure 400 {string} string "invalid request body or invalid product id"
// @Failure 401 {string} string "invalid or expired token"
// @Failure 403 {string} string "permission denied"
// @Failure 404 {string} string "data not found"
// @Failure 409 {string} string "insufficient stock"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/inventory/{id} [patch]
func (h *InventoryHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")
	productId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req model.InventoryUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	log := &model.InventoryLog{
		UserId:         claims.UserId,
		Operation:      req.Operation,
		ChangeQuantity: req.ChangeQuantity,
	}

	stockQuantity, err := h.service.UpdateStock(r.Context(), productId, log)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Message       string `json:"message"`
		StockQuantity int    `json:"stock_quantity"`
	}{
		Message:       "Inventory updated successfully",
		StockQuantity: stockQuantity,
	}); err != nil {
		writeError(w, err)
	}
}
