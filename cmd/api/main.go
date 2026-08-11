package main

import (
	"os"

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

	port := os.Getenv("PORT")
	DB := db.StartDB()
	authRepo := repository.NewAuthRepository(DB)
	authService := service.NewAuthService(authRepo)

	postRepo := repository.NewPostRepository(DB)
	postService := service.NewPostService(postRepo)

	commentRepo := repository.NewCommentRepository(DB)
	commentService := service.NewCommentService(commentRepo)

	reactionRepo := repository.NewReactionRepository(DB)
	reactionService := service.NewReactionService(reactionRepo)

	h := router.NewHandlers(authService, postService, commentService, reactionService)
	r := router.New(h, handler.AuthMiddleware())
	r.Run(":" + port)
}
