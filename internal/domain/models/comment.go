package models

type Comment struct {
	ID              string
	Author          User
	ParentCommentID *string
	Content         string
	Reactions       map[ReactionType]int
	PostID          string
	CreatedAt       uint64
}
