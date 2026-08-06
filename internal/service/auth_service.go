package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/dto"
	"github.com/alexistamher/social-api-go/internal/repository"
)

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	_, erro := s.repo.Register(dto.DtoAuthResponseToDomain(&req))
	if erro != nil {
		return dto.AuthResponse{}, erro
	}

	// TODO: generate token and return it
	return dto.AuthResponse{}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	return dto.AuthResponse{}, nil
}

func (s *authService) GetInfo(ctx context.Context, userID string) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
