package models

import "time"

type User struct {
	ID          string
	Username    string
	Email       string
	Password    string
	DisplayName string
	Bio         *string
	AvatarURL   *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Author struct {
	ID          string
	Username    string
	DisplayName string
	AvatarURL   *string
}
