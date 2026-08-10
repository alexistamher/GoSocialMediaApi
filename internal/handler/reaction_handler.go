package handler

import (
	errors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/gin-gonic/gin"
)

type ReactionHandler struct {
	reactionService service.ReactionService
}

func NewReactionHandler() *ReactionHandler {
	return &ReactionHandler{}
}

func (h *ReactionHandler) AddReaction(c *gin.Context) {
	var req dto.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.reactionService.AddReaction(c, req); err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *ReactionHandler) DeleteReaction(c *gin.Context) {
	reactionID := c.Param("reaction_id")
	if reactionID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.reactionService.DeleteReaction(c, reactionID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *ReactionHandler) UpdateReaction(c *gin.Context) {
	var req dto.UpdateReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.reactionService.UpdateReaction(c, req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *ReactionHandler) GetTargetReactions(c *gin.Context) {
	var req dto.GetTargetReactionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.reactionService.GetTargetReactions(c, req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *ReactionHandler) AddTargetReaction(c *gin.Context) {
	var req dto.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}
	err := h.reactionService.AddReaction(c, req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}
