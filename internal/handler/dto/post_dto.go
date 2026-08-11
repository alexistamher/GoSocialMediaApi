package dto

import "github.com/alexistamher/social-api-go/internal/domain/models"

type CreatePostRequest struct {
	Content    string  `json:"content" binding:"required,min=1,max=2000"`
	ParentID   *string `json:"parent_id"`
	Visibility string  `json:"visibility" binding:"required,oneof=friends public"`
	UserID     string
}

type CreatedPostResponse struct {
	ID        string `json:"id"`
	CreatedAt uint64 `json:"created_at"`
}

type PostResponse struct {
	ID               string         `json:"id"`
	Content          string         `json:"content"`
	Author           AuthorResponse `json:"author"`
	CommentsCount    int            `json:"comments_count"`
	PreviewReactions map[string]int `json:"preview_reactions"`
	Visibility       string         `json:"visibility" binding:"oneof=friends public"`
	CreatedAt        uint64         `json:"created_at"`
}

type PostWithDetailsResponse struct {
	ID         string                   `json:"id"`
	Content    string                   `json:"content"`
	Author     AuthorResponse           `json:"author"`
	Reactions  []TargetReactionResponse `json:"reactions"`
	Visibility string                   `json:"visibility" binding:"oneof=friends public"`
	CreatedAt  uint64                   `json:"created_at"`
}

func (r *CreatePostRequest) ToDomainPost() *models.Post {
	var postParent *models.Post
	if r.ParentID != nil {
		postParent = &models.Post{ID: *r.ParentID}
	}

	return &models.Post{
		Content:    r.Content,
		Visibility: models.PostVisibility(r.Visibility),
		Author:     models.Author{ID: r.UserID},
		PostParent: postParent,
	}
}

func ResponseFromDomainPost(p *models.Post) *PostResponse {
	return &PostResponse{
		ID:               p.ID,
		Content:          p.Content,
		Author:           *ResponseFromDomainAuthor(&p.Author),
		CommentsCount:    p.CommentsCount,
		PreviewReactions: p.PreviewReactions,
		Visibility:       string(p.Visibility),
		CreatedAt:        p.CreatedAt,
	}
}

func ResponseFromDomainPostWithDetails(p *models.PostWithDetails) *PostWithDetailsResponse {
	reactions := make([]TargetReactionResponse, len(p.Reactions))
	for i, reaction := range p.Reactions {
		reactions[i] = *FromDomainTargetReaction(&reaction)
	}
	return &PostWithDetailsResponse{
		ID:         p.ID,
		Content:    p.Content,
		Author:     *ResponseFromDomainAuthor(&p.Author),
		Reactions:  reactions,
		Visibility: string(p.Visibility),
		CreatedAt:  p.CreatedAt,
	}
}
