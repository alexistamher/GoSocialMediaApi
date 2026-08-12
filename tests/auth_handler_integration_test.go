package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexistamher/social-api-go/internal/handler"
	"github.com/alexistamher/social-api-go/internal/handler/dto"
	"github.com/alexistamher/social-api-go/internal/repository"
	"github.com/alexistamher/social-api-go/internal/router"
	"github.com/alexistamher/social-api-go/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Integration_Flow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	db, cleanup := setupAuthIntegrationDB(t)
	defer cleanup()

	authRepo := repository.NewAuthRepository(db)
	authSvc := service.NewAuthService(authRepo)

	handlers := router.Handlers{
		Auth: handler.NewAuthHandler(authSvc),
	}

	r := router.New(handlers, handler.AuthMiddleware())

	t.Run("1. Register User", func(t *testing.T) {
		regReq := dto.RegisterRequest{
			Username:    "integrationUser",
			Email:       "integration@example.com",
			Password:    "password123",
			DisplayName: "Integration User",
		}

		body, _ := json.Marshal(regReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var res dto.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
	})

	var accessToken string

	t.Run("2. Login User", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Email:    "integration@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res dto.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)

		accessToken = res.AccessToken
	})

	t.Run("3. Login Invalid Password - Error", func(t *testing.T) {
		loginReq := dto.LoginRequest{
			Email:    "integration@example.com",
			Password: "wrong_password",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("4. Get Protected Info with Valid JWT", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/info", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var userRes dto.UserResponse
		err := json.Unmarshal(w.Body.Bytes(), &userRes)
		assert.NoError(t, err)
		assert.Equal(t, "integrationUser", userRes.Username)
		assert.Equal(t, "Integration User", userRes.DisplayName)
	})

	t.Run("6. Get Protected Info without JWT - Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/info", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
