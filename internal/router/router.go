package router

import (
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/middleware"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

func NewRouter(cafeHandler *handler.CafeHandler, userHandler *handler.UserHandler, authenHandler *handler.AuthHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/cafes", middleware.Authenticate(cafeHandler.GetAll))
	mux.HandleFunc("GET /api/v1/cafes/{id}", cafeHandler.GetByID)
	mux.HandleFunc("POST /api/v1/cafes", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/cafes/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/cafes/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Delete)))

	mux.HandleFunc("GET /api/v1/users", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.GetAll)))
	mux.HandleFunc("GET /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.GetByID)))
	mux.HandleFunc("POST /api/v1/users", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Delete)))

	mux.HandleFunc("POST /api/v1/auth/login", authenHandler.Login)

	return mux
}
