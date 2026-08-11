package repository

import "github.com/alexistamher/social-api-go/internal/domain/models"

type ReactionRepository interface {
	AddReaction(postID string, userID string, reactionType string, reactionTargetType string) (*models.Reaction, error)
	UpdateReaction(reactionID string, reactionType string) error
	DeleteReaction(reactionID string) error
	GetTargetReactions(targetID string) ([]*models.Reaction, error)
	GetTargetPreviewReactions(targetID string) (map[string]int, error)
}
