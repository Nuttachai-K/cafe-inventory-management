package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/database"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/repository"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading configuration from the environment")
	}

	email := os.Getenv("ADMIN_EMAIL")
	password := os.Getenv("ADMIN_PASSWORD")

	username := os.Getenv("ADMIN_USERNAME")
	if username == "" {
		username = "admin"
	}

	if email == "" || password == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD must both be set")
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	userService := service.NewUserService(repository.NewUserRepository(db))

	admin := &model.User{
		Username: username,
		Email:    email,
		Password: password,
		UserRole: model.RoleAdmin,
	}

	switch err := userService.Create(ctx, admin); {
	case err == nil:
		log.Printf("created admin user %s", email)
	case errors.Is(err, service.ErrDuplicateEmail):
		log.Printf("admin user %s already exists, leaving it unchanged", email)
	default:
		log.Fatalf("create admin user: %v", err)
	}
}
