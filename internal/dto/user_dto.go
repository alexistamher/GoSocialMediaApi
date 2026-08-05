package dto

type UserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	CreatedAt   uint   `json:"created_at"`
}

type UpdateProfileRequest struct {
	DisplayName *string `json:"display_name" binding:"omitempty,min=1,max=64"`
	Bio         *string `json:"bio" binding:"omitempty,max=280"`
	AvatarURL   *string `json:"avatar_url" binding:"omitempty,url"`
}

type SearchUsersQuery struct {
	Query  string `form:"q" binding:"required,min=1"`
	Limit  int    `form:"limit,default=20" binding:"min=1,max=50"`
	Offset int    `form:"offset,default=0" binding:"min=0"`
}
