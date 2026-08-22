package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GetInfo(ctx context.Context, userID string) (*dto.UserResponse, error)
}

type CommentService interface {
	AddComment(c *gin.Context, req dto.AddCommentRequest) error
	GetCommentByID(c *gin.Context, commentID string) error
	GetCommentsByCommentID(c *gin.Context, commentID string) error
	GetCommentsByPostID(c *gin.Context, postID string) error
	DeleteComment(c *gin.Context, commentID string) error
}
