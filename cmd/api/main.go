package main

import (
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/router"
	"github.com/alexistamher/social-api-go/internal/service"
)

func main() {
	authService := service.NewAuthService()
	h := router.NewHandlers(authService)
	r := router.New(h, handler.AuthMiddleware())
	r.Run(":3000")
}
