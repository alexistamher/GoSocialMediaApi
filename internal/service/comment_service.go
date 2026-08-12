package service

import (
	"net/http"

	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

func (s *commentService) AddComment(c *gin.Context, req dto.AddCommentRequest) error {
	result, err := s.repo.AddComment(req.ToDomainComment())
	if err != nil {
		return err
	}
	reponseBody := dto.ResponseFromDomainComment(result)
	c.JSON(http.StatusOK, reponseBody)
	return nil
}

func (s *commentService) GetCommentByID(c *gin.Context, commentID string) error {
	result, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}
	responseBody := dto.ResponseFromDomainCommentWitDetails(result)
	c.JSON(http.StatusOK, responseBody)
	return nil
}

func (s *commentService) DeleteComment(c *gin.Context, commentID string) error {
	err := s.repo.DeleteComment(commentID)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	return nil
}

func (s *commentService) GetCommentsByPostID(c *gin.Context, postID string) error {
	comments, err := s.repo.GetCommentsByPostID(postID)
	if err != nil {
		return err
	}

	rcomments := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		rcomments[i] = *dto.ResponseFromDomainComment(comment)
	}
	c.JSON(http.StatusOK, rcomments)
	return nil
}
