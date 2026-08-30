package service

import (
	"context"
	"log"
	"net/http"
	"time"

	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/domain/notification"
	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

const userIDContextKey = "userID"

const notifyTimeout = 3 * time.Second

type ReactionService struct {
	repo     repository.ReactionRepository
	notifier notification.Notifier
}

func NewReactionService(repo repository.ReactionRepository, notifier notification.Notifier) ReactionService {
	if notifier == nil {
		notifier = notification.NoopNotifier{}
	}
	return ReactionService{
		repo:     repo,
		notifier: notifier,
	}
}

func (r *ReactionService) AddReaction(ctx *gin.Context, req dto.AddReactionRequest) error {
	res, err := r.repo.AddReaction(req.TargetID, req.UserID, req.ReactionType, req.ReactionTargetType)
	if err != nil {
		return err
	}
	r.notifyReaction(ctx, res, notification.ActionAdd)
	responseBody := dto.FromDomainTargetReaction(res)
	ctx.JSON(http.StatusOK, responseBody)
	return nil
}

func (r *ReactionService) DeleteReaction(ctx *gin.Context, reactionID string) error {
	reaction, err := r.repo.DeleteReaction(reactionID)
	if err != nil {
		return err
	}
	r.notifyReaction(ctx, reaction, notification.ActionDelete)
	ctx.JSON(http.StatusOK, gin.H{"message": "Reaction deleted successfully"})
	return nil
}

func (r *ReactionService) UpdateReaction(ctx *gin.Context, req dto.UpdateReactionRequest) error {
	reaction, err := r.repo.UpdateReaction(req.ID, req.ReactionType)
	if err != nil {
		return err
	}
	r.notifyReaction(ctx, reaction, notification.ActionUpdate)
	ctx.JSON(http.StatusOK, gin.H{"message": "Reaction updated successfully"})
	return nil
}

func (r *ReactionService) GetTargetReactions(ctx *gin.Context, req dto.GetTargetReactionsRequest) error {
	reactions, err := r.repo.GetTargetReactions(req.TargetID)
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

func (r *ReactionService) notifyReaction(ctx *gin.Context, reaction *dmodels.Reaction, action notification.Action) {
	if reaction == nil {
		return
	}
	connectionID, _ := ctx.Value(userIDContextKey).(string)
	evt := notification.ReactionEvent{
		ConnectionID:   connectionID,
		ReactionID:     reaction.ID,
		ParentTargetID: reaction.TargetID,
		Action:         action,
	}

	go func() {
		bg, cancel := context.WithTimeout(context.Background(), notifyTimeout)
		defer cancel()
		if err := r.notifier.NotifyReaction(bg, evt); err != nil {
			log.Printf("reaction notification failed (reaction=%s action=%s): %v", evt.ReactionID, action, err)
		}
	}()
}
