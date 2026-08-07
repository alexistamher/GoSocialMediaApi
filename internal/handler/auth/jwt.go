package auth

import (
	"os"
	"time"

	errors "github.com/alexistamher/social-api-go/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"user_id"`
	Jti    string `json:"jti"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	claims := Claims{
		UserID: userID,
		Jti:    uuid.NewString(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.ErrUnauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.ErrUnauthorized
	}
	return claims, nil
}
