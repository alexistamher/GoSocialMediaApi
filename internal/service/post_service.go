package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/handler/dto"
)

type PostService struct {
}

func (p PostService) DeletePost(context context.Context, postID string) error {
	panic("unimplemented")
}

func (p PostService) GetUserPosts(context context.Context, userID string) ([]*dto.PostResponse, error) {
	panic("unimplemented")
}

func (p PostService) CreatePost(context context.Context, req dto.CreatePostRequest) (*dto.CreatedCommentResponse, error) {
	panic("unimplemented")
}

func (p PostService) GetPostByID(context context.Context, postID string) (*dto.PostResponse, error) {
	panic("unimplemented")
}
