package repository

import "github.com/alexistamher/social-api-go/internal/domain/models"

type ReactionRepository interface {
	AddReaction(req *models.Reaction) (*models.Reaction, error)
	DeleteReaction(reactionID string) error
	UpdateReaction(req *models.Reaction) error
	GetTargetReactions(targetID string, targetType string) ([]*models.Reaction, error)
}
