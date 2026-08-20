package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Nuttachai-K/cafe-inventory-management/docs"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/middleware"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
)

func NewRouter(serverHandler *handler.ServerHandler, cafeHandler *handler.CafeHandler, userHandler *handler.UserHandler, authenHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler, categoryHandler *handler.CategoryHandler, inventoryHandler *handler.InventoryHandler,
	inventoryLogHandler *handler.InventoryLogHandler) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", serverHandler.Health)
	mux.HandleFunc("GET /readyz", serverHandler.Ready)

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

	mux.HandleFunc("GET /api/v1/inventory/{id}", inventoryHandler.GetByID)
	mux.HandleFunc("GET /api/v1/inventory", inventoryHandler.GetAll)
	mux.HandleFunc("PATCH /api/v1/inventory/{id}", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(inventoryHandler.UpdateStock)))

	mux.HandleFunc("GET /api/v1/inventory/logs", middleware.Authenticate(middleware.RequireRole(model.RoleAdmin)(inventoryLogHandler.GetLogs)))

	mux.HandleFunc("POST /api/v1/auth/login", authenHandler.Login)

	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
