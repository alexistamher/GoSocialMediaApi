package dto

import dmodels "github.com/alexistamher/social-api-go/internal/domain/models"

type AddCommentRequest struct {
	Content         string  `json:"content" binding:"required,min=1,max=1000"`
	ParentCommentID *string `json:"parent_comment_id,omitempty" binding:"omitempty,uuid"`
	PostID          string  `json:"post_id" binding:"required,uuid"`
	UserID          string
}

type CreatedCommentResponse struct {
	ID string `json:"id"`
}

type CommentResponse struct {
	ID               string                    `json:"id"`
	Content          string                    `json:"content"`
	Author           AuthorResponse            `json:"author"`
	CommentsCount    uint                      `json:"comments_count"`
	PreviewReactions []PreviewReactionResponse `json:"preview_reactions"`
	CreatedAt        uint64                    `json:"created_at"`
	PostID           string                    `json:"post_id"`
	ParentCommentID  *string                   `json:"parent_comment_id"`
}

type CommentDetailsResponse struct {
	ID        string                   `json:"id"`
	Reactions []TargetReactionResponse `json:"reactions"`
	Comments  []CommentResponse        `json:"comments"`
}

type GetCommentsResponse struct {
	Comments []CommentResponse `json:"comments"`
}

func (c *AddCommentRequest) ToDomainComment() *dmodels.Comment {
	return &dmodels.Comment{
		Content:         c.Content,
		ParentCommentID: c.ParentCommentID,
		PostID:          c.PostID,
		Author:          dmodels.Author{ID: c.UserID},
	}
}

func ResponseFromDomainComment(c *dmodels.Comment) *CommentResponse {
	preactions := make([]PreviewReactionResponse, len(c.PreviewReactions))
	for i, reaction := range c.PreviewReactions {
		preactions[i] = *ResponseFromDomainPreviewReaction(&reaction)
	}

	return &CommentResponse{
		ID:               c.ID,
		Content:          c.Content,
		Author:           *ResponseFromDomainAuthor(&c.Author),
		CommentsCount:    c.CommentsCount,
		PreviewReactions: preactions,
		CreatedAt:        c.CreatedAt,
		PostID:           c.PostID,
		ParentCommentID:  c.ParentCommentID,
	}
}

func ResponseFromDomainCommentWitDetails(c *dmodels.CommentWithDetails) *CommentDetailsResponse {
	reactions := make([]TargetReactionResponse, len(c.Reactions))
	comments := make([]CommentResponse, len(c.Comments))
	for i, react := range c.Reactions {
		reactions[i] = *FromDomainTargetReaction(&react)
	}

	return &CommentDetailsResponse{
		ID:        c.ID,
		Reactions: reactions,
		Comments:  comments,
	}
}
