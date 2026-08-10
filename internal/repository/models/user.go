package models

import (
	"time"

	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Username     string    `gorm:"column:username;not null"`
	Email        string    `gorm:"column:email;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	DisplayName  string    `gorm:"column:display_name;not null"`
	Bio          *string   `gorm:"column:bio"`
	AvatarURL    *string   `gorm:"column:avatar_url"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null;default:now()"`
}

type Authors struct {
	ID          string `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Username    string `gorm:"column:username;not null"`
	DisplayName string `gorm:"column:display_name;not null"`
	AvatarURL   string `gorm:"column:avatar_url"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}

func DomainUserToEntity(u *dmodels.User) *User {
	return &User{
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.Password,
		DisplayName:  u.DisplayName,
		Bio:          u.Bio,
		AvatarURL:    u.AvatarURL,
	}
}

func (u *User) ToDomainUser() *dmodels.User {
	return &dmodels.User{
		ID:          u.ID.String(),
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		AvatarURL:   u.AvatarURL,
	}
}
