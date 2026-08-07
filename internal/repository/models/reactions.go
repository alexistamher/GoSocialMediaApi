package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reactions struct {
	gorm.Model

	ID                 uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ReactionType       string    `gorm:"column:reaction_type;not null"`
	ReactionTargetType string    `gorm:"column:reaction_target_type;not null"`
	UserID             uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	TargetID           uuid.UUID `gorm:"column:target_id;type:uuid;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()"`
}
