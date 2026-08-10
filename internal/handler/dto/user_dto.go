package dto

import "github.com/alexistamher/social-api-go/internal/domain/models"

type UserResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	Bio         *string `json:"bio,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	CreatedAt   uint    `json:"created_at"`
	UpdatedAt   uint    `json:"updated_at"`
}

func UserDomainToDto(user *models.User) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		CreatedAt:   uint(user.CreatedAt.Unix()),
		UpdatedAt:   uint(user.UpdatedAt.Unix()),
	}
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,min=1,max=64"`
	Bio         *string `json:"bio" binding:"omitempty,max=280"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
}

type SearchUsersQuery struct {
	Query  string `form:"q" binding:"required,min=1"`
	Limit  int    `form:"limit,default=20" binding:"min=1,max=50"`
	Offset int    `form:"offset,default=0" binding:"min=0"`
}

type AuthorResponse struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

func ResponseFromDomainAuthor(a *models.Author) *AuthorResponse {
	return &AuthorResponse{
		ID:          a.ID,
		Username:    a.Username,
		DisplayName: a.DisplayName,
		AvatarURL:   a.AvatarURL,
	}
}
