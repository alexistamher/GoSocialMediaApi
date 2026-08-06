package repository

import (
	"github.com/alexistamher/social-api-go/internal/domain/models"
)

type AuthRepository interface {
	Register(user *models.User) (*string, error)
	Login(email string) (*models.Token, error)
	GetUserInfo(userId string) (*models.User, error)
}
