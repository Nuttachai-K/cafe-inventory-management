package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/database"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/handler"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/router"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/utils"
)

func main() {
	// Load enviorment data
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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

	mux := router.NewRouter(cafeHandler, userHandler, authHandler, productHandler, categoryHandler, inventoryHandler)

	http.ListenAndServe(":8080", mux)
}
