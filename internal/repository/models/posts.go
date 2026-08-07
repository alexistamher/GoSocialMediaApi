package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Posts struct {
	gorm.Model

	ID           uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	Content      string     `gorm:"column:content;not null"`
	Visibility   string     `gorm:"column:visibility;not null"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null;default:now()"`
	AuthorID     uuid.UUID  `gorm:"column:author_id;type:uuid;not null"`
	PostParentID *uuid.UUID `gorm:"column:post_parent_id;type:uuid"`
}
