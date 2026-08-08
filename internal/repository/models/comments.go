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

func EntityFromCommentDomain(comment *dmodels.Comment) *Comments {
	var parentCommentID *uuid.UUID
	if comment.ParentCommentID != nil {
		ID := uuid.MustParse(*comment.ParentCommentID)
		parentCommentID = &ID
	}

	return &Comments{
		Content:         comment.Content,
		AuthorID:        uuid.MustParse(comment.Author.ID),
		PostID:          uuid.MustParse(comment.PostID),
		ParentCommentID: parentCommentID,
		CreatedAt:       time.Unix(int64(comment.CreatedAt), 0),
	}
}

func (c *Comments) ToDomainComment() *dmodels.Comment {
	var parentCommentID *string
	if c.ParentCommentID != nil {
		ID := c.ParentCommentID.String()
		parentCommentID = &ID
	}

	return &dmodels.Comment{
		ID:              c.ID.String(),
		Content:         c.Content,
		Author:          dmodels.User{ID: c.AuthorID.String()},
		PostID:          c.PostID.String(),
		ParentCommentID: parentCommentID,
		CreatedAt:       uint64(c.CreatedAt.Unix()),
	}
}
