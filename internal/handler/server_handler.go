package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerHandler struct {
	db *pgxpool.Pool
}

func NewServerHandler(db *pgxpool.Pool) *ServerHandler {
	return &ServerHandler{
		db: db,
	}
}

// Health godoc
// @Summary Liveness check
// @Tags server
// @Success 200
// @Router /health [get]
func (h *ServerHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Ready godoc
// @Summary Readiness check
// @Tags server
// @Success 200
// @Failure 503
// @Router /readyz [get]
func (h *ServerHandler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}
