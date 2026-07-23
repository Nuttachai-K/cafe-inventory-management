package router

import (
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
)

func NewRouter(cafeHandler *handler.CafeHandler, userHandler *handler.UserHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cafes", cafeHandler.GetAll)
	mux.HandleFunc("GET /cafes/{id}", cafeHandler.GetByID)
	mux.HandleFunc("POST /cafes", cafeHandler.Create)
	mux.HandleFunc("PATCH /cafes/{id}", cafeHandler.Update)
	mux.HandleFunc("DELETE /cafes/{id}", cafeHandler.Delete)

	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("GET /users/{id}", userHandler.GetByID)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("PATCH /users/{id}", userHandler.Update)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)

	return mux
}
