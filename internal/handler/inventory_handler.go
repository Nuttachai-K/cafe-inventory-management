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
