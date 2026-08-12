package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	derrors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/handler/auth"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/handler/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("JWT_SECRET", "test-secret-key"); err != nil {
		panic(err)
	}
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func validToken(t *testing.T) (token string, userID string) {
	t.Helper()
	userID = "ae122adc-22c0-4d4b-a3c2-170ef99cfb5c"
	var err error
	token, err = auth.GenerateToken(userID)
	if err != nil {
		t.Fatalf("failed to generate test token: %v", err)
	}
	return token, userID
}

func setupAuthRouter(h *handler.AuthHandler) *gin.Engine {
	r := gin.New()

	public := r.Group("/auth")
	public.POST("/register", h.Register)
	public.POST("/login", h.Login)

	private := r.Group("/auth")
	private.Use(handler.AuthMiddleware())
	private.GET("/info", h.GetInfo)

	return r
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
	expected := &dto.AuthResponse{
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
	assert.Equal(t, *expected, got)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Register_ValidationError_InvalidEmail(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	invalidBody := []byte(`{"username":"jconnor","email":"not-an-email","password":"supersecret","display_name":"John Connor"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
}

func TestAuthHandler_Register_ValidationError_ShortPassword(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	invalidBody := []byte(`{"username":"jconnor","email":"jconnor@example.com","password":"123","display_name":"John Connor"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
}

func TestAuthHandler_Register_ValidationError_MissingFields(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	invalidBody := []byte(`{"email":"jconnor@example.com","password":"supersecret"}`)
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

	svc.On("Register", mock.Anything, reqBody).Return((*dto.AuthResponse)(nil), derrors.ErrAlreadyExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), derrors.ErrAlreadyExists.Error())
	svc.AssertExpectations(t)
}

func TestAuthHandler_Register_ServiceError_Internal(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.RegisterRequest{
		Username:    "john",
		Email:       "jconnor@example.com",
		Password:    "supersecret",
		DisplayName: "John Connor",
	}

	svc.On("Register", mock.Anything, reqBody).Return((*dto.AuthResponse)(nil), assert.AnError)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.LoginRequest{
		Email:    "jconnor@example.com",
		Password: "supersecret",
	}
	expected := &dto.AuthResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
	}

	svc.On("Login", mock.Anything, reqBody).Return(expected, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got dto.AuthResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, *expected, got)
	svc.AssertExpectations(t)
}

func TestAuthHandler_Login_ValidationError_InvalidEmail(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	invalidBody := []byte(`{"email":"not-an-email","password":"supersecret"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(invalidBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	svc.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	reqBody := dto.LoginRequest{
		Email:    "jconnor@example.com",
		Password: "wrong-password",
	}

	svc.On("Login", mock.Anything, reqBody).Return((*dto.AuthResponse)(nil), derrors.ErrInvalidCredentials)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), derrors.ErrInvalidCredentials.Error())
	svc.AssertExpectations(t)
}

func TestAuthHandler_GetInfo_Success(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	token, userID := validToken(t)

	bio := "I'll be back"
	avatarURL := "https://example.com/avatar.jpg"
	expected := &dto.UserResponse{
		ID:          userID,
		Username:    "john",
		DisplayName: "John Connor",
		Bio:         &bio,
		AvatarURL:   &avatarURL,
	}

	svc.On("GetInfo", mock.Anything, userID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got dto.UserResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, *expected, got)
	svc.AssertExpectations(t)
}

func TestAuthHandler_GetInfo_NoAuthorizationHeader(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetInfo", mock.Anything, mock.Anything)
}

func TestAuthHandler_GetInfo_InvalidTokenFormat(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	req.Header.Set("Authorization", "invalid-token-without-bearer-prefix")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token format")
	svc.AssertNotCalled(t, "GetInfo", mock.Anything, mock.Anything)
}

func TestAuthHandler_GetInfo_InvalidJWTSignature(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.invalid-signature")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	svc.AssertNotCalled(t, "GetInfo", mock.Anything, mock.Anything)
}

func TestAuthHandler_GetInfo_ServiceError_NotFound(t *testing.T) {
	svc := new(mocks.AuthServiceMock)
	h := handler.NewAuthHandler(svc)
	router := setupAuthRouter(h)

	token, userID := validToken(t)

	svc.On("GetInfo", mock.Anything, userID).Return((*dto.UserResponse)(nil), derrors.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/auth/info", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), derrors.ErrNotFound.Error())
	svc.AssertExpectations(t)
}
