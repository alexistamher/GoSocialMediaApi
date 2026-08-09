package dto

import "github.com/google/uuid"

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
	ID           uuid.UUID    `json:"id"`
	TargetID     uuid.UUID    `json:"target_id"`
	ReactionType string       `json:"reaction_type"`
	CreatedAt    uint64       `json:"created_at"`
	Author       AuthResponse `json:"author"`
}

type UpdatedReactionResponse struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt uint64    `json:"updated_at"`
}

type TargetPreviewReactionResponse struct {
	Reactions map[string]int `json:"reactions"`
}
