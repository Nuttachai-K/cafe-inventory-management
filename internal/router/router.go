package router

import (
	"net/http"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/middleware"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

func NewRouter(cafeHandler *handler.CafeHandler, userHandler *handler.UserHandler, authenHandler *handler.AuthHandler, productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/cafes", cafeHandler.GetAll)
	mux.HandleFunc("GET /api/v1/cafes/{id}", cafeHandler.GetByID)
	mux.HandleFunc("POST /api/v1/cafes", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/cafes/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/cafes/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(cafeHandler.Delete)))

	mux.HandleFunc("GET /api/v1/users", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.GetAll)))
	mux.HandleFunc("GET /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.GetByID)))
	mux.HandleFunc("POST /api/v1/users", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/users/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(userHandler.Delete)))

	mux.HandleFunc("GET /api/v1/products", productHandler.GetAll)
	mux.HandleFunc("GET /api/v1/products/{id}", productHandler.GetByID)
	mux.HandleFunc("POST /api/v1/products", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(productHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/products/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(productHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/products/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(productHandler.Delete)))

	mux.HandleFunc("GET /api/v1/categories", categoryHandler.GetAll)
	mux.HandleFunc("GET /api/v1/categories/{id}", categoryHandler.GetByID)
	mux.HandleFunc("POST /api/v1/categories", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(categoryHandler.Create)))
	mux.HandleFunc("PATCH /api/v1/categories/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(categoryHandler.Update)))
	mux.HandleFunc("DELETE /api/v1/categories/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(categoryHandler.Delete)))

	mux.HandleFunc("POST /api/v1/auth/login", authenHandler.Login)

	return mux
}
