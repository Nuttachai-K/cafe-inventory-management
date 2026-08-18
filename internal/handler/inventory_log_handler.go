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

// GetLogs godoc
// @Summary Get all logs
// @Tags log
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param inventory_id query string false "Inventory id"
// @Param user_id query string false "User id"
// @Param operation query string false "Operation"
// @Param from query string false "Starting time"
// @Param to query string false "End time"
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.InventoryLog
// @Failure 400 {string} string "invalid query parameter or value"
// @Failure 401 {string} string "invalid or expired token"
// @Failure 403 {string} string "permission denied"
// @Failure 500 {string} string "internal server error"
// @Router /api/v1/inventory/logs [get]
func (h *InventoryLogHandler) GetLogs(w http.ResponseWriter, r *http.Request) {

	var filter model.InventoryLogFilter

	if v := r.URL.Query().Get("inventory_id"); v != "" {
		id, err := strconv.Atoi(v)
		if err != nil {
			writeJSONError(
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
			writeJSONError(
				w,
				"Invalid user id",
				http.StatusBadRequest,
			)
			return
		}
		filter.UserId = &id
	}

	if v := r.URL.Query().Get("operation"); v != "" {
		filter.Operation = &v
	}

	if v := r.URL.Query().Get("from"); v != "" {
		from, err := time.Parse(time.DateOnly, v) // "2006-01-02"
		if err != nil {
			writeJSONError(w, "invalid from date", http.StatusBadRequest)
			return
		}
		filter.From = &from
	}

	if v := r.URL.Query().Get("to"); v != "" {
		to, err := time.Parse(time.DateOnly, v)
		if err != nil {
			writeJSONError(w, "invalid to date", http.StatusBadRequest)
			return
		}
		filter.To = &to
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			writeJSONError(w, "invalid limit", http.StatusBadRequest)
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
