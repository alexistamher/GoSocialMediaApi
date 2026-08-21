package dto

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/google/uuid"
)

type AddReactionRequest struct {
	TargetID           string `json:"target_id" binding:"required,uuid"`
	ReactionType       string `json:"reaction_type" binding:"required,oneof=like love haha wow sad angry"`
	ReactionTargetType string `json:"reaction_target_type" binding:"required,oneof=post comment"`
	UserID             string
}

type UpdateReactionRequest struct {
	ID           string `json:"id" binding:"required,uuid"`
	ReactionType string `json:"reaction_type" binding:"required,oneof=like love haha wow sad angry"`
}

type CreatedReactionResponse struct {
	ID        string `json:"id"`
	CreatedAt uint64 `json:"created_at"`
}

type TargetReactionResponse struct {
	ID           string         `json:"id"`
	TargetID     string         `json:"target_id"`
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

type PreviewReactionResponse struct {
	ID           string `json:"id"`
	ReactionType string `json:"reaction_type"`
	TargetID     string `json:"target_id"`
	AuthorID     string `json:"author_id"`
}

func (r *AddReactionRequest) ToDomainReaction() *models.Reaction {
	return &models.Reaction{
		TargetID:           r.TargetID,
		ReactionType:       models.ReactionType(r.ReactionType),
		ReactionTargetType: models.ReactionTargetType(r.ReactionTargetType),
	}
}

func (r *UpdateReactionRequest) ToDomainReaction() *models.Reaction {
	return &models.Reaction{
		ID:           r.ID,
		ReactionType: models.ReactionType(r.ReactionType),
	}
}

func FromDomainTargetReaction(r *models.Reaction) *TargetReactionResponse {
	return &TargetReactionResponse{
		ID:           r.ID,
		TargetID:     r.TargetID,
		ReactionType: string(r.ReactionType),
		CreatedAt:    r.CreatedAt,
		Author:       AuthorResponse(r.Author),
	}
}

func ResponseFromDomainPreviewReaction(r *models.PreviewReaction) *PreviewReactionResponse {
	return &PreviewReactionResponse{
		ID:           r.ID,
		ReactionType: string(r.ReactionType),
		TargetID:     r.TargetID,
		AuthorID:     r.AuthorID,
	}
}
