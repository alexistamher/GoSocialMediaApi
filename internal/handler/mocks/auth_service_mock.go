package mocks

import (
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(ctx *gin.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *AuthServiceMock) Login(ctx *gin.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*dto.AuthResponse), args.Error(1)
}

func (m *AuthServiceMock) GetInfo(ctx *gin.Context, userID string) (*dto.UserResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *AuthServiceMock) Health(ctx *gin.Context, userID string) {}
