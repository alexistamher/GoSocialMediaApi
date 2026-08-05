package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/dto"
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/handler/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAuthRouter(h *handler.AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(handler.AuthMiddleware())
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/auth/info", h.GetInfo)
	return r
}

func TestAuthHandler_GetInfo_FailGettingUserId(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	svc.On("GetInfo", mock.Anything, "").Return(dto.UserResponse{}, domain.ErrMissingUserID)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Contains(t, w.Body.String(), domain.ErrMissingUserID.Error())
	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertExpectations(t)
}

func TestAuthHandler_GetInfo_Success(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	expected := dto.UserResponse{
		ID:          "user-123",
		Username:    "john",
		DisplayName: "John Connor",
		Bio:         "Bio",
		AvatarURL:   "https://example.com/avatar.jpg",
		CreatedAt:   uint(time.Now().Unix()),
	}

	svc.On("GetInfo", mock.Anything, "1").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	req.Header.Set("Authorization", "1")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var got dto.UserResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, expected, got)
	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Register_Success(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.RegisterRequest{
		Username:    "john",
		Email:       "jconnor@example.com",
		Password:    "supersecret",
		DisplayName: "John Connor",
	}
	expected := dto.AuthResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	svc.On("Register", mock.Anything, reqBody).Return(expected, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var got dto.AuthResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, expected, got)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationError(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	invalidBody := []byte(`{"username":"jconnor","email":"bad-email","password":"123","display_name":"John Connor"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
}

func TestAuthHandler_Register_ServiceError_AlreadyExists(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.RegisterRequest{
		Username:    "john",
		Email:       "jconnor@example.com",
		Password:    "supersecret",
		DisplayName: "John Connor",
	}

	svc.On("Register", mock.Anything, reqBody).Return(dto.AuthResponse{}, domain.ErrAlreadyExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.LoginRequest{
		Email:    "jconnor@example.com",
		Password: "wrong-password",
	}

	svc.On("Login", mock.Anything, reqBody).Return(dto.AuthResponse{}, domain.ErrInvalidCredentials)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertExpectations(t)
}
