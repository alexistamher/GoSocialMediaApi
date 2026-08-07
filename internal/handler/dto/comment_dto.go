package dto

type CreateCommentRequest struct {
	Content         string  `json:"content" binding:"required,min=1,max=1000"`
	ParentCommentID *string `json:"parent_comment_id,omitempty" binding:"omitempty,uuid"`
}

type CommentResponse struct {
	ID               string         `json:"id"`
	Content          string         `json:"content"`
	Author           CommentAuthor  `json:"author"`
	ReactionsPreview map[string]int `json:"reactions"`
	CreatedAt        uint64         `json:"created_at"`
	PostID           string         `json:"post_id"`
	ParentCommentID  *string        `json:"parent_comment_id" binding:"uuid"`
}

type CommentAuthor struct {
	ID          string  `json:"id"`
	DisplayName string  `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}
