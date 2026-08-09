package handler

import (
	"net/http"

	errors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{commentService: commentService}
}

func (h *CommentHandler) AddComment(c *gin.Context) {
	var req dto.AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, errors.ErrBadRequest)
		return
	}

	comment, err := h.commentService.AddComment(req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	if commentID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.commentService.DeleteComment(commentID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	comments, err := h.commentService.GetPostComments(postID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}
