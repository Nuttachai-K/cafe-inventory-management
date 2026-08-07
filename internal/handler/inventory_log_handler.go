package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type InventoryLogHandler struct {
	service service.InventoryLogService
}

func NewInventoryLogHandler(service service.InventoryLogService) *InventoryLogHandler {
	return &InventoryLogHandler{
		service: service,
	}
}

func (h *InventoryLogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {

	var filter model.InventoryLogFilter

	if v := r.URL.Query().Get("inventory_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			http.Error(
				w,
				"Invalid inventory id",
				http.StatusBadRequest,
			)
			return
		}
		filter.InventoryId = &id
	}

	if v := r.URL.Query().Get("user_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			http.Error(
				w,
				"Invalid user id",
				http.StatusBadRequest,
			)
			return
		}
		filter.UserId = &id
	}

	if v := r.URL.Query().Get("from"); v != "" {
		from, err := time.Parse(time.DateOnly, v) // "2006-01-02"
		if err != nil {
			http.Error(w, "invalid from date", http.StatusBadRequest)
			return
		}
		filter.From = &from
	}

	if v := r.URL.Query().Get("to"); v != "" {
		to, err := time.Parse(time.DateOnly, v)
		if err != nil {
			http.Error(w, "invalid to date", http.StatusBadRequest)
			return
		}
		filter.To = &to
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	logs, err := h.service.GetLog(r.Context(), &filter, limit)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		writeError(w, err)
	}
}
