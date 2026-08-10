package repository

import "github.com/alexistamher/social-api-go/internal/domain/models"

type CommentRepository interface {
	AddComment(comment *models.Comment) (*models.Comment, error)
	DeleteComment(commentID string) error
	GetCommentsByPostID(postID string) ([]models.Comment, error)
}
