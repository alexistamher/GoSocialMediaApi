package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(user *dmodels.User) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	userEtt := models.DomainUserToEntity(user)

	ans := r.db.Create(&userEtt)
	if ans.Error != nil {
		return nil, ans.Error
	}

	id := userEtt.ID.String()
	return &id, nil
}

func (r *authRepository) Login(email string) (*dmodels.Token, error) {
	return nil, nil
}

func (r *authRepository) GetUserInfo(userId string) (*dmodels.User, error) {
	return nil, nil
}
