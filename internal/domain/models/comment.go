package models

type Comment struct {
	ID              string
	Author          Author
	ParentCommentID *string
	Content         string
	Reactions       map[string]int
	PostID          string
	CreatedAt       uint64
}

type CommentWithAuthor struct {
	ID               string
	Author           *Author
	ParentCommentID  *string
	Content          string
	PreviewReactions map[string]int
	PostID           string
	CreatedAt        uint64
}
