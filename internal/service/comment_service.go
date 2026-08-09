package service

import (
	"github.com/alexistamher/social-api-go/internal/handler/dto"
)

type CommentService struct {
}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (c *CommentService) AddComment(req dto.AddCommentRequest) (*dto.CreatedCommentResponse, error) {
	panic("unimplemented")
}

func (c *CommentService) GetPostComments(postID string) ([]*dto.CommentResponse, error) {
	panic("unimplemented")
}

func (c *CommentService) DeleteComment(commentID string) error {
	panic("unimplemented")
}
