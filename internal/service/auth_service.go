package service

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/auth"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
)

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	userID, erro := s.repo.Register(dto.DtoAuthResponseToDomain(&req))
	if erro != nil {
		return nil, erro
	}
	token, erro := auth.GenerateToken(*userID)
	if erro != nil {
		return nil, erro
	}

	return &dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: token,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	userID, erro := s.repo.Login(req.Email, req.Password)
	if erro != nil {
		return nil, erro
	}
	token, erro := auth.GenerateToken(*userID)
	if erro != nil {
		return nil, erro
	}

	return &dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: token,
	}, nil
}

func (s *authService) GetInfo(ctx context.Context, userID string) (*dto.UserResponse, error) {
	user, erro := s.repo.GetUserInfo(userID)
	if erro != nil {
		return nil, erro
	}

	return dto.UserDomainToDto(user), nil
}
