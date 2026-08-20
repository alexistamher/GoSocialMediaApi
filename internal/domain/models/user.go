package models

type User struct {
	ID          string
	Username    string
	Email       string
	Password    string
	DisplayName string
	Bio         *string
	AvatarURL   *string
	CreatedAt   uint64
	UpdatedAt   uint64
}

type Author struct {
	ID          string
	Username    string
	DisplayName string
	AvatarURL   *string
}
