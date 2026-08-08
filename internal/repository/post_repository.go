package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	drepository "github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) drepository.PostRepository {
	return &postRepository{db: db}
}

func (p *postRepository) AddPost(post *dmodels.Post) (*dmodels.Post, error) {
	postEntity := models.EntityFromPostDomain(post)
	if err := p.db.Create(&postEntity).Error; err != nil {
		return nil, err
	}
	return postEntity.ToDomainPost(), nil
}

func (p *postRepository) DeletePost(postID string) error {
	if err := p.db.Where("id = ?", postID).Delete(&models.Posts{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetAllPosts(userID string, offset *int, limit *int) ([]*dmodels.Post, *int, error) {
	var posts []*models.Posts
	var total int64

	baseQuery := p.db.Model(&models.Posts{}).Where("author_id = ?", userID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	query := baseQuery.Order("created_at desc")
	if offset != nil {
		query = query.Offset(*offset)
	}
	if limit != nil {
		query = query.Limit(*limit)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, nil, err
	}

	dpost := make([]*dmodels.Post, len(posts))
	for i, post := range posts {
		dpost[i] = post.ToDomainPost()
	}

	currentOffset := 0
	if offset != nil {
		currentOffset = *offset
	}
	nextOffset := currentOffset + len(posts)

	return dpost, &nextOffset, nil
}

func (p *postRepository) GetPostByID(postID string) (*dmodels.Post, error) {
	var post *models.Posts
	if err := p.db.Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	return post.ToDomainPost(), nil
}

func (p *postRepository) GetPostsByUserID(userID string, offset *int, limit uint) ([]*dmodels.Post, *int, error) {
	var posts []*models.Posts
	var total int64

	baseQuery := p.db.Model(&models.Posts{}).Where("author_id = ?", userID)

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	query := baseQuery.Order("created_at desc")
	if offset != nil {
		query = query.Offset(*offset)
	}
	query = query.Limit(int(limit))

	if err := query.Find(&posts).Error; err != nil {
		return nil, nil, err
	}

	dpost := make([]*dmodels.Post, len(posts))
	for i, post := range posts {
		dpost[i] = post.ToDomainPost()
	}

	currentOffset := 0
	if offset != nil {
		currentOffset = *offset
	}
	nextOffset := currentOffset + len(posts)

	return dpost, &nextOffset, nil
}

func (p *postRepository) AddReaction(postID string, userID string, reactionType string, reactionTargetType string) (*dmodels.Reaction, error) {
	var reaction = models.Reactions{
		TargetID:           uuid.MustParse(postID),
		UserID:             uuid.MustParse(userID),
		ReactionType:       reactionType,
		ReactionTargetType: reactionTargetType,
	}

	if err := p.db.Create(&reaction).Error; err != nil {
		return nil, err
	}
	return models.EntityFromReactionDomain(&reaction), nil
}

func (p *postRepository) DeleteReaction(reactionID string) error {
	var reaction models.Reactions
	if err := p.db.Where("id = ?", reactionID).Delete(&reaction).Error; err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetTargetReactions(targetID string) ([]*dmodels.Reaction, error) {
	var reactions []models.Reactions
	if err := p.db.Where("target_id = ?", targetID).Find(&reactions).Error; err != nil {
		return nil, err
	}
	dreactions := make([]*dmodels.Reaction, len(reactions))
	for i, reaction := range reactions {
		dreactions[i] = models.EntityFromReactionDomain(&reaction)
	}
	return dreactions, nil
}

func (p *postRepository) GetTargetPreviewReactions(targetID string) (map[string]int, error) {
	var reactions []models.Reactions
	if err := p.db.Where("target_id = ?", targetID).Find(&reactions).Error; err != nil {
		return nil, err
	}
	reactionsPreview := make(map[string]int)
	for _, reaction := range reactions {
		reactionsPreview[reaction.ReactionType]++
	}
	return reactionsPreview, nil
}

func (p *postRepository) AddComment(comment *dmodels.Comment) (*dmodels.Comment, error) {
	commentEntity := models.EntityFromCommentDomain(comment)
	if err := p.db.Create(commentEntity).Error; err != nil {
		return nil, err
	}
	return commentEntity.ToDomainComment(), nil
}

func (p *postRepository) DeleteComment(commentID string) error {
	var children *[]models.Comments
	if err := p.db.Where("parent_comment_id = ?", commentID).Find(&children).Error; err != nil {
		return err
	}

	for _, c := range *children {
		if err := p.DeleteComment(c.ID.String()); err != nil {
			return err
		}
	}

	if err := p.db.Where("id = ?", commentID).Delete(&models.Comments{}).Error; err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetCommentsByPostID(postID string) ([]dmodels.Comment, error) {
	var comments []models.Comments
	if err := p.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		return nil, err
	}

	dcomments := make([]dmodels.Comment, len(comments))
	for i, c := range comments {
		dcomments[i] = *c.ToDomainComment()
	}

	return dcomments, nil
}
