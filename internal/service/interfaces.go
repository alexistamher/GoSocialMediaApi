package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error)
	GetInfo(ctx context.Context, userID string) (dto.UserResponse, error)
}
