package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/handler/dto"
)

type ReactionService struct {
}

func NewReactionService() *ReactionService {
	return &ReactionService{}
}

func (r *ReactionService) AddReaction(ctx context.Context, req dto.AddReactionRequest) (*dto.CreatedReactionResponse, error) {
	panic("unimplemented")
}

func (r *ReactionService) DeleteReaction(ctx context.Context, reactionID string) error {
	panic("unimplemented")
}

func (r *ReactionService) UpdateReaction(ctx context.Context, req dto.UpdateReactionRequest) error {
	panic("unimplemented")
}

func (r *ReactionService) GetTargetReactions(ctx context.Context, userID string) ([]*dto.TargetReactionResponse, error) {
	panic("unimplemented")
}

func (r *ReactionService) GetTargetPreviewReactions(ctx context.Context, postId string) (*dto.TargetPreviewReactionResponse, error) {
	panic("unimplemented")
}
