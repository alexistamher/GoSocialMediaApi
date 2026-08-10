package models

type PostVisibility string

const (
	Friends PostVisibility = "friends"
	Public  PostVisibility = "public"
)

type Post struct {
	ID               string
	Content          string
	Author           Author
	CommentsCount    int
	PreviewReactions map[string]int
	Visibility       PostVisibility
	PostParent       *Post
	CreatedAt        uint64
}
