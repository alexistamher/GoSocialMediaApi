package handler

import (
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

	req.UserID = c.Value(UserIDKey).(string)
	if err := h.commentService.AddComment(c, req); err != nil {
		respondError(c, err)
		return
	}
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("comment_id")
	if commentID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.commentService.DeleteComment(c, commentID); err != nil {
		respondError(c, err)
		return
	}
}

func (h *CommentHandler) GetCommentById(c *gin.Context) {
	commentID := c.Param("comment_id")
	if commentID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.commentService.GetCommentByID(c, commentID); err != nil {
		respondError(c, err)
		return
	}
}

func (h *CommentHandler) GetPostComments(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.commentService.GetCommentsByPostID(c, postID); err != nil {
		respondError(c, err)
		return
	}
}
