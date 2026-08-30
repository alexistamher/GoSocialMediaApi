package main

import (
	"log"
	"os"

	"github.com/alexistamher/social-api-go/internal/client/notificationgrpc"
	"github.com/alexistamher/social-api-go/internal/domain/notification"
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

	notifier := buildNotifier()
	reactionRepo := repository.NewReactionRepository(DB)
	reactionService := service.NewReactionService(reactionRepo, notifier)

	h := router.NewHandlers(authService, postService, commentService, reactionService)
	r := router.New(h, handler.AuthMiddleware())
	_ = r.Run(":" + port)
}

func buildNotifier() notification.Notifier {
	addr := os.Getenv("NOTIFICATION_GRPC_ADDR")
	if addr == "" {
		log.Println("NOTIFICATION_GRPC_ADDR not set; reaction notifications disabled")
		return notification.NoopNotifier{}
	}

	plaintext := os.Getenv("NOTIFICATION_GRPC_PLAINTEXT") == "true"
	client, err := notificationgrpc.NewClient(addr, plaintext)
	if err != nil {
		log.Printf("notification client init failed (%s); reaction notifications disabled: %v", addr, err)
		return notification.NoopNotifier{}
	}
	log.Printf("reaction notifications enabled -> %s (plaintext=%t)", addr, plaintext)
	return client
}
