package repository

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
)

type PostRepository interface {
	AddPost(post *models.Post) (*models.Post, error)
	DeletePost(postID string) error
	GetAllPosts(userID string, offset *int, limit *int) ([]*models.Post, *int, error)
	GetPostByID(postID string) (*models.Post, error)
	GetPostsByUserID(userID string, offset *int, limit *int) ([]*models.Post, *int, error)
	AddPostReaction(postID string, reactionType string) (*models.Reaction, error)
	DeletePostReaction(postID string) error
	AddComment(comment *models.Comment) (*models.Comment, error)
	DeleteComment(commentID string) error
	GetCommentsByPostID(postID string) ([]models.Comment, error)
}
