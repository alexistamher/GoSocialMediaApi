package models

type ReactionType string
type ReactionTargetType string

const (
	PostType    ReactionTargetType = "post"
	CommentType ReactionTargetType = "comment"
)

const (
	LikeType  ReactionType = "like"
	LoveType  ReactionType = "love"
	HahaType  ReactionType = "haha"
	WowType   ReactionType = "wow"
	SadType   ReactionType = "sad"
	AngryType ReactionType = "angry"
)

type Reaction struct {
	ID                 string
	ReactionType       ReactionType
	ReactionTargetType ReactionTargetType
	UserID             string
	TargetID           string
}
