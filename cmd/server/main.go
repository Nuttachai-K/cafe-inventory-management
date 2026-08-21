package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/Nuttachai-K/cafe-inventory-management/docs"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/database"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/middleware"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/router"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

// @title Cafe Inventory Management API
// @version 1.0
// @description REST API for multi-branch cafe chain inventory
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load enviorment data
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, reading configuration from the environment")
	}

	err = utils.CheckJWTSecret()
	if err != nil {
		log.Fatal(err)
	}

	// Connect with database
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if host := os.Getenv("SWAGGER_HOST"); host != "" {
		docs.SwaggerInfo.Host = host
		scheme := os.Getenv("SWAGGER_SCHEME")
		if scheme == "" {
			scheme = "https"
		}
		docs.SwaggerInfo.Schemes = []string{scheme}
	}

	cafeRepo := repository.NewCafeRepository(db)
	cafeService := service.NewCafeService(cafeRepo)
	cafeHandler := handler.NewCafeHandler(cafeService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo, db)
	productHandler := handler.NewProductHandler(productService)

	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo, db)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	inventoryLogRepo := repository.NewInventoryLogRepository(db)
	inventoryLogService := service.NewInventoryLogService(inventoryLogRepo)
	inventoryLogHandler := handler.NewInventoryLogHandler(inventoryLogService)

	serverHandler := handler.NewServerHandler(db)

	mux := router.NewRouter(serverHandler, cafeHandler, userHandler, authHandler, productHandler, categoryHandler, inventoryHandler, inventoryLogHandler)

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		log.Fatal("SERVER_ADDRESS environment variable is not set")
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: middleware.CorsMiddleware(mux),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log.Printf("starting server on %s", srv.Addr)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received, draining connections")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
