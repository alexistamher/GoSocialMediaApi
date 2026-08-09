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

	res, err := h.postService.CreatePost(c.Request.Context(), req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *PostHandler) GetPostByID(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}
	res, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("post_id")
	if postID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}

	err := h.postService.DeletePost(c.Request.Context(), postID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID := c.Value("user_id").(string)
	if userID == "" {
		respondError(c, errors.ErrBadRequest)
		return
	}
	posts, err := h.postService.GetUserPosts(c.Request.Context(), userID)
	if err != nil {
		respondError(c, errors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, posts)
}
