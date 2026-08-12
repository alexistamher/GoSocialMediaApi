package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	errors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/alexistamher/social-api-go/internal/handler/auth"
	"github.com/gin-gonic/gin"
)

type contextKey = string

const UserIDKey contextKey = "userID"
const TokenIdKey contextKey = "tokenIdKey"

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			err, _ := json.Marshal(map[string]any{
				"error": "invalid token format",
			})
			ctx.Writer.Header().Add("Content-Type", "application/json")
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			_, _ = ctx.Writer.Write(err)
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.ErrInvalidCredentials)
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(UserIDKey, claims.UserID)
		ctx.Set(TokenIdKey, TokenId{claims.Jti, claims.ExpiresAt.Time})
		ctx.Next()
	}
}

type TokenId struct {
	Jti       string
	ExpiresAt time.Time
}
