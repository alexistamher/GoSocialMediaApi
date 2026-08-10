package service

import (
	"net/http"

	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type CommentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

func (s *CommentService) AddComment(c *gin.Context, req dto.AddCommentRequest) error {
	result, err := s.repo.AddComment(req.ToDomainComment())
	if err != nil {
		return err
	}
	reponseBody := dto.ResponseFromDomainComment(result)
	c.JSON(http.StatusOK, reponseBody)
	return nil
}

func (s *CommentService) GetCommentsByCommentID(c *gin.Context, commentID string) error {
	result, err := s.repo.GetCommentsByCommentID(commentID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, result)
	return nil
}

func (s *CommentService) DeleteComment(c *gin.Context, commentID string) error {
	err := s.repo.DeleteComment(commentID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	return nil
}
