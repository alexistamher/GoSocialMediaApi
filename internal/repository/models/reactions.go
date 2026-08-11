package models

import (
	"time"

	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/google/uuid"
)

type Reactions struct {
	ID                 uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid();not null"`
	ReactionType       string    `gorm:"column:reaction_type;not null"`
	ReactionTargetType string    `gorm:"column:target_type;not null"`
	UserID             uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	TargetID           uuid.UUID `gorm:"column:target_id;type:uuid;not null"`
	CreatedAt          time.Time `gorm:"column:created_at;not null;default:now()"`
	DeletedAt          time.Time `gorm:"column:deleted_at"`

	Author User `gorm:"foreignKey:UserID;references:ID"`
}

func EntityFromReactionDomain(reaction *Reactions) *dmodels.Reaction {
	return &dmodels.Reaction{
		ID:                 reaction.ID.String(),
		ReactionType:       dmodels.ReactionType(reaction.ReactionType),
		ReactionTargetType: dmodels.ReactionTargetType(reaction.ReactionTargetType),
		UserID:             reaction.UserID.String(),
		TargetID:           reaction.TargetID.String(),
		CreatedAt:          uint64(reaction.CreatedAt.UnixMilli()),
		Author:             *reaction.Author.ToDomainAuthor(),
	}
}

func (r *Reactions) ToDomainReaction() *dmodels.Reaction {
	return &dmodels.Reaction{
		ID:                 r.ID.String(),
		ReactionType:       dmodels.ReactionType(r.ReactionType),
		ReactionTargetType: dmodels.ReactionTargetType(r.ReactionTargetType),
		UserID:             r.UserID.String(),
		TargetID:           r.TargetID.String(),
		CreatedAt:          uint64(r.CreatedAt.UnixMilli()),
		Author:             *r.Author.ToDomainAuthor(),
	}
}
