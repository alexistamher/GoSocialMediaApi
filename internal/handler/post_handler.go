package handler

import (
	"net/http"

	errors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = c.Value(UserIDKey).(string)
	if err := h.postService.CreatePost(c, req); err != nil {
		respondError(c, err)
		return
	}
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.postService.GetPostByID(c, postID); err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.postService.DeletePost(c, postID); err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := c.Value("user_id").(string)
	if userID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.postService.GetUserPosts(c, userID); err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
}
