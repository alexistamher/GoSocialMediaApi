package dto

import "github.com/alexistamher/social-api-go/internal/domain/models"

type RegisterRequest struct {
	Username    string  `json:"username" binding:"required,min=3,max=32,alphanum"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	DisplayName string  `json:"display_name" binding:"required,min=1,max=64"`
	Bio         *string `json:"bio" binding:"omitempty,max=500"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func DtoAuthResponseToDomain(authReq *RegisterRequest) *models.User {
	return &models.User{
		Username:    authReq.Username,
		Password:    authReq.Password,
		Email:       authReq.Email,
		DisplayName: authReq.DisplayName,
		Bio:         authReq.Bio,
		AvatarURL:   authReq.AvatarURL,
	}
}
