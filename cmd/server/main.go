package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/Nuttachai-K/cafe-finder-api/internal/database"
	"github.com/Nuttachai-K/cafe-finder-api/internal/handler"
	"github.com/Nuttachai-K/cafe-finder-api/internal/repository"
	"github.com/Nuttachai-K/cafe-finder-api/internal/router"
	"github.com/Nuttachai-K/cafe-finder-api/internal/service"
)

func main() {
	// Load enviorment data
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect with database
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewCafeRepository(db)
	service := service.NewCafeService(repo)
	handler := handler.NewCafeHandler(service)

	mux := router.NewRouter(handler)

	http.ListenAndServe(":8080", mux)
}
