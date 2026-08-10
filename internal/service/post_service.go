package service

import (
	"net/http"

	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) DeletePost(c *gin.Context, postID string) error {
	if err := s.repo.DeletePost(postID); err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	return nil
}

func (s *PostService) GetUserPosts(c *gin.Context, userID string) error {
	posts, nextCursor, err := s.repo.GetPostsByUserID(userID, nil, 10)
	if err != nil {
		return err
	}

	rposts := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		rposts[i] = dto.ResponseFromDomainPost(post)
	}

	response := map[string]any{
		"next_cursor": nextCursor,
		"posts":       rposts,
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (s *PostService) CreatePost(c *gin.Context, req dto.CreatePostRequest) error {
	post := req.ToDomainPost()
	post, err := s.repo.AddPost(post)
	if err != nil {
		return err
	}

	responseBody := dto.ResponseFromDomainPost(post)
	c.JSON(http.StatusOK, responseBody)
	return nil
}

func (s *PostService) GetPostByID(c *gin.Context, postID string) error {
	post, err := s.repo.GetPostByID(postID)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, dto.ResponseFromDomainPost(post))
	return nil
}
