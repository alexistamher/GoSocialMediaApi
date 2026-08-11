package repository

import (
	dmodels "github.com/alexistamher/social-api-go/internal/domain/models"
	drepository "github.com/alexistamher/social-api-go/internal/domain/repository"
	"github.com/alexistamher/social-api-go/internal/repository/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) drepository.ReactionRepository {
	return &reactionRepository{db: db}
}

func (p *reactionRepository) AddReaction(postID string, userID string, reactionType string, reactionTargetType string) (*dmodels.Reaction, error) {
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

func (p *reactionRepository) UpdateReaction(reactionID string, reactionType string) error {
	var reaction models.Reactions
	if err := p.db.Where("id = ?", reactionID).First(&reaction).Error; err != nil {
		return err
	}
	reaction.ReactionType = reactionType
	if err := p.db.Save(&reaction).Error; err != nil {
		return err
	}
	return nil
}

func (p *reactionRepository) DeleteReaction(reactionID string) error {
	var reaction models.Reactions
	if err := p.db.Where("id = ?", reactionID).Delete(&reaction).Error; err != nil {
		return err
	}
	return nil
}

func (p *reactionRepository) GetTargetReactions(targetID string) ([]*dmodels.Reaction, error) {
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

func getTargetPreviewReactions(tx *gorm.DB, targetID string) (map[string]int, error) {
	var reactions []models.Reactions
	if err := tx.
		Where("target_id = ?", targetID).
		Select("reaction_type", "target_type").Find(&reactions).Error; err != nil {
		return nil, err
	}
	reactionsPreview := make(map[string]int)
	for _, reaction := range reactions {
		reactionsPreview[reaction.ReactionType]++
	}
	return reactionsPreview, nil
}

func (p *reactionRepository) GetTargetPreviewReactions(targetID string) (map[string]int, error) {
	return getTargetPreviewReactions(p.db, targetID)
}

func getReactionsByTargetId(db *gorm.DB, targetID string) ([]*models.Reactions, error) {
	var reactions []*models.Reactions
	if err := db.Preload("Author").Where("target_id = ?", targetID).Find(&reactions).Error; err != nil {
		return nil, err
	}

	return reactions, nil
}

func getPreviewReactionsByIDs(db *gorm.DB, targetIDs []string) (map[string]map[string]int, error) {
	reactionx := make(map[string]map[string]int)
	var reactions []models.Reactions
	if err := db.Where("target_id IN ?", targetIDs).Find(&reactions).Error; err != nil {
		return nil, err
	}

	for _, reaction := range reactions {
		if _, ok := reactionx[reaction.TargetID.String()]; !ok {
			reactionx[reaction.TargetID.String()] = map[string]int{
				reaction.ReactionType: 1,
			}
		} else {
			reactionx[reaction.TargetID.String()][reaction.ReactionType]++
		}
	}
	return reactionx, nil
}
