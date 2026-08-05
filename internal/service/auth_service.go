package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/dto"
)

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	return dto.AuthResponse{}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	return dto.AuthResponse{}, nil
}

func (s *authService) GetInfo(ctx context.Context, userID string) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
