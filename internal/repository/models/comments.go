package models

import (
	"time"

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
