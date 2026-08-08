package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
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
	if err := p.db.Where("id = ?", postID).Delete(&dmodels.Post{}).Error; err != nil {
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
	// TODO: debería ser llenado con un mapeo
	var dpost *dmodels.Post
	return dpost, nil
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

func (p *postRepository) AddPostReaction(postID string, reactionType string) (*dmodels.Reaction, error) {
	// TODO: deberia construirse a partir de un maper
	var reaction models.Reactions

	if err := p.db.Create(&reaction).Error; err != nil {
		return nil, err
	}
	// TODO: deberia maprear el reaction entity a domain
	var dReaction dmodels.Reaction
	return &dReaction, nil
}

func (p *postRepository) DeletePostReaction(postID string) error {
	var reaction models.Reactions
	if err := p.db.Where("id = ?", postID).Delete(&reaction).Error; err != nil {
		return err
	}
	return nil
}

func (p *postRepository) AddComment(comment *dmodels.Comment) (*dmodels.Comment, error) {
	commentEntity := models.EntityFromCommentDomain(comment)
	if err := p.db.Create(commentEntity).Error; err != nil {
		return nil, err
	}
	return commentEntity.ToDomainComment(), nil
}

func (p *postRepository) DeleteComment(commentID string) error {
	if err := p.db.Where("id = ?", commentID).Delete(&dmodels.Comment{}).Error; err != nil {
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
