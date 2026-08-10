package dto

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/google/uuid"
)

type AddReactionRequest struct {
	TargetID           uuid.UUID `json:"target_id" binding:"required,uuid"`
	ReactionType       string    `json:"reaction_type" binding:"required,oneof=like love haha wow sad angry"`
	ReactionTargetType string    `json:"reaction_target_type" binding:"required,oneof=post comment"`
}

type UpdateReactionRequest struct {
	ID           uuid.UUID `json:"id" binding:"required,uuid"`
	ReactionType string    `json:"reaction_type" binding:"required"`
}

type CreatedReactionResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt uint64    `json:"created_at"`
}

type TargetReactionResponse struct {
	ID           uuid.UUID      `json:"id"`
	TargetID     uuid.UUID      `json:"target_id"`
	ReactionType string         `json:"reaction_type"`
	CreatedAt    uint64         `json:"created_at"`
	Author       AuthorResponse `json:"author"`
}

type UpdatedReactionResponse struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt uint64    `json:"updated_at"`
}

type TargetPreviewReactionResponse struct {
	Reactions map[string]int `json:"reactions"`
}

type GetTargetReactionsRequest struct {
	TargetID           string `json:"target_id" binding:"required,uuid"`
	ReactionTargetType string `json:"reaction_target_type" binding:"required"`
}

func (r *AddReactionRequest) ToDomainReaction() *models.Reaction {
	return &models.Reaction{
		TargetID:           r.TargetID.String(),
		ReactionType:       models.ReactionType(r.ReactionType),
		ReactionTargetType: models.ReactionTargetType(r.ReactionTargetType),
	}
}

func (r *UpdateReactionRequest) ToDomainReaction() *models.Reaction {
	return &models.Reaction{
		ID:           r.ID.String(),
		ReactionType: models.ReactionType(r.ReactionType),
	}
}

func FromDomainReaction(r *models.Reaction) *CreatedReactionResponse {
	return &CreatedReactionResponse{
		ID:        uuid.MustParse(r.ID),
		CreatedAt: r.CreatedAt,
	}
}

func FromDomainTargetReaction(r *models.Reaction) *TargetReactionResponse {
	return &TargetReactionResponse{
		ID:           uuid.MustParse(r.ID),
		TargetID:     uuid.MustParse(r.TargetID),
		ReactionType: string(r.ReactionType),
		CreatedAt:    r.CreatedAt,
		Author:       AuthorResponse{ID: r.UserID},
	}
}
