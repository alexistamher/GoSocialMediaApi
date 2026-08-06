package repository

import (
	"log"

	errors "github.com/alexistamher/social-api-go/internal/domain"
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	drepository "github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) drepository.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(user *dmodels.User) (*string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	userEtt := models.DomainUserToEntity(user)
	userEtt.PasswordHash = string(hashedPassword)

	ans := r.db.Create(&userEtt)
	if ans.Error != nil {
		return nil, ans.Error
	}

	id := userEtt.ID.String()
	return &id, nil
}

func (r *authRepository) Login(email string, password string) (*string, error) {
	var user models.User
	if ans := r.db.Where("email = ?", email).First(&user); ans.Error != nil {
		log.Printf("error finding user in: %s", ans.Error)
		return nil, errors.ErrInvalidCredentials
	}

	log.Printf("password: %s", password)
	log.Printf("user: %v", user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("error logging in: %s", err)
		return nil, errors.ErrInvalidCredentials
	}

	id := user.ID.String()
	return &id, nil
}

func (r *authRepository) GetUserInfo(userId string) (*dmodels.User, error) {
	var userEtt models.User
	if r.db.Where("id = ?", userId).First(&userEtt).Error != nil {
		return nil, errors.ErrInvalidCredentials
	}
	return userEtt.ToDomainUser(), nil
}
