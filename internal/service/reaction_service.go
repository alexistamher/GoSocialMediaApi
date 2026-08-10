package service

import (
	"net/http"

	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type ReactionService struct {
	repo repository.ReactionRepository
}

func NewReactionService() *ReactionService {
	return &ReactionService{}
}

func (r *ReactionService) AddReaction(ctx *gin.Context, req dto.AddReactionRequest) error {
	res, err := r.repo.AddReaction(req.ToDomainReaction())
	if err != nil {
		return err
	}
	responseBody := dto.FromDomainReaction(res)
	ctx.JSON(http.StatusOK, responseBody)
	return nil
}

func (r *ReactionService) DeleteReaction(ctx *gin.Context, reactionID string) error {
	if err := r.repo.DeleteReaction(reactionID); err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Reaction deleted successfully"})
	return nil
}

func (r *ReactionService) UpdateReaction(ctx *gin.Context, req dto.UpdateReactionRequest) error {
	if err := r.repo.UpdateReaction(req.ToDomainReaction()); err != nil {
		return err
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Reaction updated successfully"})
	return nil
}

func (r *ReactionService) GetTargetReactions(ctx *gin.Context, req dto.GetTargetReactionsRequest) error {
	reactions, err := r.repo.GetTargetReactions(req.TargetID, req.ReactionTargetType)
	if err != nil {
		return err
	}
	responseBody := make([]*dto.TargetReactionResponse, len(reactions))
	for i, reaction := range reactions {
		responseBody[i] = dto.FromDomainTargetReaction(reaction)
	}
	ctx.JSON(http.StatusOK, responseBody)
	return nil
}
