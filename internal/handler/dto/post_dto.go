package dto

type CreatePostRequest struct {
	Content    string  `json:"content" binding:"required,min=1,max=2000"`
	ParentID   *string `json:"parent_id"`
	Visibility string  `json:"visibility" binding:"required,oneof=friends public"`
}

type CreatedPostResponse struct {
	ID        string `json:"id"`
	CreatedAt uint64 `json:"created_at"`
}

type PostResponse struct {
	ID               string         `json:"id"`
	Content          string         `json:"content"`
	Author           AuthResponse   `json:"author"`
	CommentsCount    int            `json:"comments_count"`
	ReactionsPreview map[string]int `json:"reactions"`
	Visibility       string         `json:"visibility" binding:"oneof=friends public"`
	CreatedAt        uint64         `json:"created_at"`
	PostParentID     *string        `json:"post_parent_id" binding:"uuid"`
}
