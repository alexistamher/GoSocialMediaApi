package main

import (
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/alexistamher/social-api-go/internal/repository/db"
	"github.com/alexistamher/social-api-go/internal/router"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	DB := db.StartDB()
	authRepo := repository.NewAuthRepository(DB)
	authService := service.NewAuthService(authRepo)
	h := router.NewHandlers(authService)
	r := router.New(h, handler.AuthMiddleware())
	r.Run(":3000")
}
