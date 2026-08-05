package mocks

import (
	"context"

	"github.com/alexistamher/social-api-go/internal/dto"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.AuthResponse), args.Error(1)
}

func (m *AuthServiceMock) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.AuthResponse), args.Error(1)
}

func (m *AuthServiceMock) GetInfo(ctx context.Context, userID string) (dto.UserResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dto.UserResponse), args.Error(1)
}
