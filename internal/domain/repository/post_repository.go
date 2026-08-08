package repository

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
)

type PostRepository interface {
	AddPost(post *models.Post) (*models.Post, error)
	DeletePost(postID string) error
	GetAllPosts(userID string, offset *int, limit *int) ([]*models.Post, *int, error)
	GetPostByID(postID string) (*models.Post, error)
	GetPostsByUserID(userID string, offset *int, limit uint) ([]*models.Post, *int, error)
	AddReaction(postID string, userID string, reactionType string, reactionTargetType string) (*models.Reaction, error)
	DeleteReaction(reactionID string) error
	GetTargetReactions(targetID string) ([]*models.Reaction, error)
	GetTargetPreviewReactions(targetID string) (map[string]int, error)
	AddComment(comment *models.Comment) (*models.Comment, error)
	DeleteComment(commentID string) error
	GetCommentsByPostID(postID string) ([]models.Comment, error)
}
