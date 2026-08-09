package dto

type AddCommentRequest struct {
	Content         string  `json:"content" binding:"required,min=1,max=1000"`
	ParentCommentID *string `json:"parent_comment_id,omitempty" binding:"omitempty,uuid"`
}

type CreatedCommentResponse struct {
	ID string `json:"id"`
}

type CommentResponse struct {
	ID               string         `json:"id"`
	Content          string         `json:"content"`
	Author           AuthResponse   `json:"author"`
	ReactionsPreview map[string]int `json:"reactions"`
	CreatedAt        uint64         `json:"created_at"`
	PostID           string         `json:"post_id"`
	ParentCommentID  *string        `json:"parent_comment_id"`
}
