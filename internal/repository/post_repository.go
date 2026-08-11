package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	drepository "github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"

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
	if err := p.db.Preload("Author").Create(&postEntity).Error; err != nil {
		return nil, err
	}
	return postEntity.ToDomainPost(), nil
}

func (p *postRepository) DeletePost(postID string) error {
	if err := p.db.Transaction(func(tx *gorm.DB) error {
		var commentIds []string
		if err := tx.Model(&models.Comments{}).
			Where("post_id = ?", postID).
			Pluck("id", &commentIds).Error; err != nil {
			return err
		}

		if err := tx.Where("target_id = ? AND target_type = 'post'", postID).
			Delete(&models.Reactions{}).Error; err != nil {
			return err
		}

		if len(commentIds) > 0 {
			for _, id := range commentIds {
				if err := deleteCommentAndChildren(tx, id); err != nil {
					return err
				}
			}
		}

		if err := tx.Where("id = ?", postID).
			Delete(&models.Posts{}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetAllPosts(userID string, offset *int, limit *int) ([]*dmodels.Post, *int, error) {
	var postIDs []string
	var posts []*models.Posts
	// var total int64

	if err := p.db.Table("posts").
		Select("id").
		Where("author_id = ?", userID).
		Pluck("id", &postIDs).Error; err != nil {
		return nil, nil, err
	}

	baseQuery := p.db.Preload("Author").Where("id IN ?", postIDs)

	// if err := baseQuery.Table("posts").Count(&total).Error; err != nil {
	// 	return nil, nil, err
	// }

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

	reactions, err := getPreviewReactionsByIDs(p.db, postIDs)
	if err != nil {
		return nil, nil, err
	}

	dposts := make([]*dmodels.Post, len(posts))
	for i, post := range posts {
		dpost := post.ToDomainPost()
		dpost.PreviewReactions = reactions[post.ID.String()]
		dposts[i] = dpost
	}

	currentOffset := 0
	if offset != nil {
		currentOffset = *offset
	}
	nextOffset := currentOffset + len(posts)

	return dposts, &nextOffset, nil
}

func (p *postRepository) GetPostByID(postID string) (*dmodels.PostWithDetails, error) {
	var post *models.Posts
	if err := p.db.Preload("Author").Where("id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}
	reactions, err := getReactionsByTargetId(p.db, postID)
	if err != nil {
		return nil, err
	}

	rpost := post.ToDomainPostWithDetails(reactions)

	return rpost, nil
}

func (p *postRepository) GetPostsByUserID(userID string, offset *int, limit uint) ([]*dmodels.Post, *int, error) {
	var postIDs []string
	var posts []*models.Posts
	var total int64

	if err := p.db.Table("posts").Where("author_id = ?", userID).Pluck("id", &postIDs).Error; err != nil {
		return nil, nil, err
	}

	baseQuery := p.db.Preload("Author").Model(&models.Posts{}).Where("id IN ?", postIDs)

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

	reactions, err := getPreviewReactionsByIDs(p.db, postIDs)
	if err != nil {
		return nil, nil, err
	}

	dpost := make([]*dmodels.Post, len(posts))
	for i, post := range posts {
		post := post.ToDomainPost()
		post.PreviewReactions = reactions[post.ID]
		dpost[i] = post
	}

	currentOffset := 0
	if offset != nil {
		currentOffset = *offset
	}
	nextOffset := currentOffset + len(posts)

	return dpost, &nextOffset, nil
}

func (p *postRepository) GetCommentsByPostID(postID string) ([]*dmodels.CommentWithAuthor, error) {
	var commentIds []string
	if err := p.db.Model(&models.Comments{}).
		Where("post_id = ? and parent_comment_id is null", postID).
		Pluck("id", &commentIds).
		Find(&models.Comments{}).Error; err != nil {
		return nil, err
	}

	var comments []models.CommentsWithAuthor
	if err := p.db.Preload("Author").Where("id in ?", commentIds).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	reactions, err := getPreviewReactionsByIDs(p.db, commentIds)
	if err != nil {
		return nil, err
	}

	dcomments := make([]*dmodels.CommentWithAuthor, len(comments))
	for i, c := range comments {
		comment := c.ToDomainCommentWithAuthor()
		comment.PreviewReactions = reactions[c.ID.String()]
		dcomments[i] = comment
	}

	return dcomments, nil
}
