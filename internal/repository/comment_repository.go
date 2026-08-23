package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"github.com/google/uuid"
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
	if err := tx.Unscoped().Model(&models.Reactions{}).
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

	if err := tx.Unscoped().Where("id = ?", commentID).Delete(&models.Comments{}).Error; err != nil {
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

	ccounters, err := getCommentsCountByCommentsIDs(p.db, commentIDs)
	if err != nil {
		return nil, err
	}

	dcomments := make([]*dmodels.Comment, len(comments))
	for i, c := range comments {
		counter := ccounters[c.ID.String()]
		comment := c.ToDomainComment()
		comment.PreviewReactions = preactions[c.ID.String()]
		comment.CommentsCount = uint(counter)
		dcomments[i] = comment
	}

	return dcomments, nil
}

func (p *commentRepository) GetCommentsByCommentID(commentID string) ([]*dmodels.Comment, error) {
	var comments []models.CommentsWithAuthor
	if err := p.db.Preload("Author").
		Raw(`
	        WITH RECURSIVE comment_tree AS (
	            SELECT id, content, author_id, post_id, parent_comment_id, created_at, 1 AS depth
	            FROM comments
	            WHERE parent_comment_id = ?

	            UNION ALL

	            SELECT c.id, c.content, c.author_id, c.post_id, c.parent_comment_id, c.created_at, ct.depth + 1
	            FROM comments c
	            INNER JOIN comment_tree ct ON c.parent_comment_id = ct.id
	        )
	        SELECT * FROM comment_tree ORDER BY depth, created_at
		`, commentID).Find(&comments).Error; err != nil {
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

	ccounters, err := getCommentsCountByCommentsIDs(p.db, commentIDs)
	if err != nil {
		return nil, err
	}

	dcomments := make([]*dmodels.Comment, len(comments))
	for i, c := range comments {
		counter := ccounters[c.ID.String()]
		comment := c.ToDomainComment()
		comment.PreviewReactions = preactions[c.ID.String()]
		comment.CommentsCount = uint(counter)
		dcomments[i] = comment
	}

	return dcomments, nil
}

func getCommentsCountByCommentsIDs(db *gorm.DB, commentIDs []string) (map[string]int, error) {
	if len(commentIDs) == 0 {
		return map[string]int{}, nil
	}

	var results []struct {
		CommentID uuid.UUID
		Count     int
	}
	if err := db.Raw(`
		WITH RECURSIVE comment_tree AS (
            SELECT id, parent_comment_id AS root_id
            FROM comments
            WHERE parent_comment_id IN ?

            UNION ALL

            SELECT c.id, ct.root_id
            FROM comments c
            INNER JOIN comment_tree ct ON c.parent_comment_id = ct.id
        )
        SELECT root_id as comment_id, COUNT(*) AS count
        FROM comment_tree
        GROUP BY root_id
	`, commentIDs).Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, r := range results {
		counts[r.CommentID.String()] = r.Count
	}
	return counts, nil
}
