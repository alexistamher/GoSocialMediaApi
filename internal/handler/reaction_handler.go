package handler

import (
	"net/http"

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

	res, err := h.reactionService.AddReaction(c.Request.Context(), req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ReactionHandler) DeleteReaction(c *gin.Context) {
	reactionID := c.Param("reaction_id")
	if reactionID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.reactionService.DeleteReaction(c.Request.Context(), reactionID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reaction deleted successfully"})
}

func (h *ReactionHandler) UpdateReaction(c *gin.Context) {
	var req dto.UpdateReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.reactionService.UpdateReaction(c.Request.Context(), req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reaction updated successfully"})
}

func (h *ReactionHandler) GetTargetReactions(c *gin.Context) {
	userID := c.Value("user_id").(string)
	reactions, err := h.reactionService.GetTargetReactions(c.Request.Context(), userID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, reactions)
}

func (h *ReactionHandler) GetTargetPreviewReactions(c *gin.Context) {
	postId := c.Param("post_id")
	if postId == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}
	reactions, err := h.reactionService.GetTargetPreviewReactions(c.Request.Context(), postId)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, reactions)
}

func (h *ReactionHandler) AddTargetReaction(c *gin.Context) {
	var req dto.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}
	rreaction, err := h.reactionService.AddReaction(c.Request.Context(), req)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, rreaction)
}
