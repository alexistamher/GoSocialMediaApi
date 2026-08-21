package models

import (
	"time"

	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comments struct {
	gorm.Model

	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Content         string     `gorm:"column:content;not null"`
	AuthorID        uuid.UUID  `gorm:"column:author_id;type:uuid;not null"`
	PostID          uuid.UUID  `gorm:"column:post_id;type:uuid;not null"`
	ParentCommentID *uuid.UUID `gorm:"column:parent_comment_id;type:uuid"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`
}

type CommentsWithAuthor struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Content         string     `gorm:"column:content;not null"`
	AuthorID        uuid.UUID  `gorm:"column:author_id;type:uuid;not null"`
	PostID          uuid.UUID  `gorm:"column:post_id;type:uuid;not null"`
	ParentCommentID *uuid.UUID `gorm:"column:parent_comment_id;type:uuid"`
	CreatedAt       time.Time  `gorm:"column:created_at;not null;default:now()"`

	Author Users `gorm:"foreignKey:AuthorID;references:ID"`
}

func (CommentsWithAuthor) TableName() string {
	return "comments"
}

func EntityFromCommentDomain(comment *dmodels.Comment) *Comments {
	var parentCommentID *uuid.UUID
	if comment.ParentCommentID != nil {
		ID, _ := uuid.Parse(*comment.ParentCommentID)
		parentCommentID = &ID
	}

	return &Comments{
		AuthorID:        uuid.MustParse(comment.Author.ID),
		PostID:          uuid.MustParse(comment.PostID),
		Content:         comment.Content,
		ParentCommentID: parentCommentID,
	}
}

func (c *Comments) ToDomainComment(user *Users) *dmodels.Comment {
	var parentCommentID *string
	var author *dmodels.Author
	if c.ParentCommentID != nil {
		ID := c.ParentCommentID.String()
		parentCommentID = &ID
	}

	if user != nil {
		author = user.ToDomainAuthor()
	} else {
		author = &dmodels.Author{ID: c.AuthorID.String()}
	}

	return &dmodels.Comment{
		ID:              c.ID.String(),
		Content:         c.Content,
		Author:          *author,
		PostID:          c.PostID.String(),
		ParentCommentID: parentCommentID,
		CreatedAt:       uint64(c.CreatedAt.Unix()),
	}
}

func (c *CommentsWithAuthor) ToDomainComment() *dmodels.Comment {
	var parentCommentID *string
	if c.ParentCommentID != nil {
		ID := c.ParentCommentID.String()
		parentCommentID = &ID
	}

	return &dmodels.Comment{
		ID:               c.ID.String(),
		PostID:           c.PostID.String(),
		ParentCommentID:  parentCommentID,
		Author:           *c.Author.ToDomainAuthor(),
		Content:          c.Content,
		PreviewReactions: []dmodels.PreviewReaction{},
		CreatedAt:        uint64(c.CreatedAt.Unix()),
	}
}
