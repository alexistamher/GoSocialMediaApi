package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) repository.CommentRepository {
	return &commentRepository{db: db}
}

func (p *commentRepository) AddComment(comment *dmodels.Comment) (*dmodels.Comment, error) {
	commentEntity := models.EntityFromCommentDomain(comment)
	if err := p.db.Create(&commentEntity).Error; err != nil {
		return nil, err
	}

	var author *models.Users
	if err := p.db.
		Select("id", "display_name", "username").
		Where("id = ?", commentEntity.AuthorID).Find(&author).Error; err != nil {
		return nil, err
	}

	return commentEntity.ToDomainComment(author), nil
}

func deleteCommentAndChildren(tx *gorm.DB, commentID string) error {
	var commentIds []string
	if err := tx.Model(&models.Comments{}).
		Where("parent_comment_id = ?", commentID).
		Pluck("id", &commentIds).Error; err != nil {
		return err
	}

	if len(commentIds) > 0 {
		for _, id := range commentIds {
			if err := deleteCommentAndChildren(tx, id); err != nil {
				return err
			}
		}
	}

	var reactionIds []string
	if err := tx.Model(&models.Reactions{}).
		Where("target_id = ? AND target_type = 'comment'", commentID).
		Pluck("id", &reactionIds).Error; err != nil {
		return err
	}

	if len(reactionIds) > 0 {
		for _, id := range reactionIds {
			if err := tx.Where("id = ?", id).Delete(&models.Reactions{}).Error; err != nil {
				return err
			}
		}
	}

	if err := tx.Where("id = ?", commentID).Delete(&models.Comments{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *commentRepository) DeleteComment(commentID string) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if err := deleteCommentAndChildren(tx, commentID); err != nil {
			return err
		}

		return nil
	})
}

func (p *commentRepository) GetCommentByID(commentID string) (*dmodels.CommentWithDetails, error) {
	var comments []models.CommentsWithAuthor
	if err := p.db.Preload("Author").
		Where("parent_comment_id = ?", commentID).Find(&comments).Error; err != nil {
		return nil, err
	}

	var reactions []models.Reactions
	if err := p.db.Preload("Author").
		Where("target_id = ?", commentID).Find(&reactions).Error; err != nil {
		return nil, err
	}

	dreactions := make([]dmodels.Reaction, len(reactions))
	for i, react := range reactions {
		dreactions[i] = *react.ToDomainReaction()
	}

	dcomments := make([]dmodels.Comment, len(comments))
	for i, c := range comments {
		dcomments[i] = *c.ToDomainComment()
	}

	cmnt := &dmodels.CommentWithDetails{
		ID:        commentID,
		Reactions: dreactions,
		Comments:  dcomments,
	}

	return cmnt, nil
}

func (p *commentRepository) GetCommentsByPostID(postID string) ([]*dmodels.Comment, error) {
	var comments []models.CommentsWithAuthor
	if err := p.db.Preload("Author").
		Where("post_id = ? AND parent_comment_id is NULL", postID).Find(&comments).Error; err != nil {
		return nil, err
	}

	commentIDs := make([]string, len(comments))
	for i, comment := range comments {
		commentIDs[i] = comment.ID.String()
	}

	preactions, err := getPreviewReactionsByIDs(p.db, commentIDs)
	if err != nil {
		return nil, err
	}

	dcomments := make([]*dmodels.Comment, len(comments))
	for i, c := range comments {
		comment := c.ToDomainComment()
		comment.PreviewReactions = preactions[c.ID.String()]
		dcomments[i] = comment
	}

	return dcomments, nil
}
