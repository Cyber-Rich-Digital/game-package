package repository

import (
	"gorm.io/gorm"
)

func NewPromotionRepository(db *gorm.DB) PromotionRepository {
	return &repo{db}
}

type PromotionRepository interface {
	GetPromotions() (interface{}, error)
}

func (r repo) GetPromotions() (interface{}, error) {

	var promotions interface{}
	if err := r.db.Find(&promotions).Error; err != nil {
		return nil, err
	}

	return promotions, nil
}
