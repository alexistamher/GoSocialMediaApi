package service

import (
	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/handler/auth"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
)

type authService struct {
	repo     repository.AuthRepository
	ntfyRepo repository.NotificationRepository
}

func NewAuthService(repo repository.AuthRepository, ntfyRepo repository.NotificationRepository) AuthService {
	return &authService{
		repo:     repo,
		ntfyRepo: ntfyRepo,
	}
}

func (s *authService) Register(ctx *gin.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
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

func (s *authService) Login(ctx *gin.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	userID, erro := s.repo.Login(req.Email, req.Password)
	if erro != nil {
		return nil, erro
	}
	token, erro := auth.GenerateToken(*userID)
	if erro != nil {
		return nil, erro
	}

	go s.ntfyRepo.RegisterConnection(ctx, *userID)
	return &dto.AuthResponse{
		AccessToken:  token,
		RefreshToken: token,
	}, nil
}

func (s *authService) GetInfo(ctx *gin.Context, userID string) (*dto.UserResponse, error) {
	user, erro := s.repo.GetUserInfo(userID)
	if erro != nil {
		return nil, erro
	}

	return dto.UserDomainToDto(user), nil
}

func (s *authService) Health(ctx *gin.Context, userID string) {
	s.ntfyRepo.RegisterConnection(ctx, userID)
}
