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
	ID               string                    `json:"id"`
	Content          string                    `json:"content"`
	Author           AuthorResponse            `json:"author"`
	CommentsCount    int                       `json:"comments_count"`
	PreviewReactions []PreviewReactionResponse `json:"preview_reactions"`
	Visibility       string                    `json:"visibility"`
	CreatedAt        uint64                    `json:"created_at"`
}

type PostWithDetailsResponse struct {
	ID               string                    `json:"id"`
	Content          string                    `json:"content"`
	Author           AuthorResponse            `json:"author"`
	PreviewReactions []PreviewReactionResponse `json:"preview_reactions"`
	Visibility       string                    `json:"visibility" binding:"oneof=friends public"`
	CreatedAt        uint64                    `json:"created_at"`
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
	rpreviewReactions := make([]PreviewReactionResponse, len(p.PreviewReactions))
	for i, preview := range p.PreviewReactions {
		rpreviewReactions[i] = *ResponseFromDomainPreviewReaction(&preview)
	}

	return &PostResponse{
		ID:               p.ID,
		Content:          p.Content,
		Author:           *ResponseFromDomainAuthor(&p.Author),
		CommentsCount:    p.CommentsCount,
		PreviewReactions: rpreviewReactions,
		Visibility:       string(p.Visibility),
		CreatedAt:        p.CreatedAt,
	}
}

func ResponseFromDomainPostWithDetails(p *models.PostWithDetails) *PostWithDetailsResponse {
	reactions := make([]PreviewReactionResponse, len(p.PreviewReactions))
	for i, reaction := range p.PreviewReactions {
		reactions[i] = *ResponseFromDomainPreviewReaction(&reaction)
	}
	return &PostWithDetailsResponse{
		ID:               p.ID,
		Content:          p.Content,
		Author:           *ResponseFromDomainAuthor(&p.Author),
		PreviewReactions: reactions,
		Visibility:       string(p.Visibility),
		CreatedAt:        p.CreatedAt,
	}
}
