package router

import (
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
)

func NewRouter(cafeHandler *handler.CafeHandler, userHandler *handler.UserHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/cafes", cafeHandler.GetAll)
	mux.HandleFunc("GET /api/v1/cafes/{id}", cafeHandler.GetByID)
	mux.HandleFunc("POST /api/v1/cafes", cafeHandler.Create)
	mux.HandleFunc("PATCH /api/v1/cafes/{id}", cafeHandler.Update)
	mux.HandleFunc("DELETE /api/v1/cafes/{id}", cafeHandler.Delete)

	mux.HandleFunc("GET /api/v1/users", userHandler.GetAll)
	mux.HandleFunc("GET /api/v1/users/{id}", userHandler.GetByID)
	mux.HandleFunc("POST /api/v1/users", userHandler.Create)
	mux.HandleFunc("PATCH /api/v1/users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", userHandler.Delete)

	return mux
}
