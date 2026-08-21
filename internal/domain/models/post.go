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
	PreviewReactions []PreviewReaction
	Visibility       PostVisibility
	PostParent       *Post
	CreatedAt        uint64
}

type PostWithDetails struct {
	ID               string
	Content          string
	Author           Author
	PreviewReactions []PreviewReaction
	Comments         []Comment
	Visibility       PostVisibility
	CreatedAt        uint64
}
