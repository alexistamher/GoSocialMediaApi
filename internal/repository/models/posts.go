package models

import (
	"time"

	"github.com/alexistamher/social-api-go/internal/domain/models"
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

	CommentChildren []*Comments `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
	Author          User        `gorm:"foreignKey:AuthorID;references:ID"`
}

func EntityFromPostDomain(p *models.Post) *Posts {
	var postParentID uuid.UUID
	if p.PostParent != nil {
		postParentID = uuid.MustParse(p.PostParent.ID)
	}
	return &Posts{
		Content:      p.Content,
		Visibility:   string(p.Visibility),
		AuthorID:     uuid.MustParse(p.Author.ID),
		PostParentID: &postParentID,
	}
}

func (p *Posts) ToDomainPost() *models.Post {
	var parent *models.Post
	if p.PostParentID != nil {
		parent = &models.Post{ID: p.PostParentID.String()}
	}
	var author models.Author
	if p.Author.ID != uuid.Nil {
		author = *p.Author.ToDomainAuthor()
	} else {
		author = models.Author{ID: p.AuthorID.String()}
	}

	return &models.Post{
		ID:               p.ID.String(),
		Content:          p.Content,
		Visibility:       models.PostVisibility(p.Visibility),
		CreatedAt:        uint64(p.CreatedAt.Unix()),
		Author:           author,
		PostParent:       parent,
		CommentsCount:    0,                // TODO: esto deberia ser obtenido
		PreviewReactions: map[string]int{}, // TODO: esto deberia ser obtenido
	}
}
