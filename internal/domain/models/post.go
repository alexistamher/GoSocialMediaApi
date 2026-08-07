package models

type PostVisibility string

const (
	Friends PostVisibility = "friends"
	Public  PostVisibility = "public"
)

type Post struct {
	ID               string
	Content          string
	Author           User
	CommentsCount    int
	ReactionsPreview map[ReactionType]int
	Visibility       PostVisibility
	PostParent       *Post
	CreatedAt        uint64
}
