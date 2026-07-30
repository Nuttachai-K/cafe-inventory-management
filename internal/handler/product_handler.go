package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProducthandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

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
