package repository

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
)

type PostRepository interface {
	AddPost(post *models.Post) (*models.Post, error)
	DeletePost(postID string) error
	GetAllPosts(userID string, offset *int, limit *int) ([]*models.Post, *int, error)
	GetPostByID(postID string) (*models.PostWithDetails, error)
	GetPostsByUserID(userID string, offset *int, limit uint) ([]*models.Post, *int, error)
	GetCommentsByPostID(postID string) ([]*models.Comment, error)
}
