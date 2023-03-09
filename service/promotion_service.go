package service

import (
	"cyber-api/model"
	"cyber-api/repository"
)

type PromotionService interface {
	GetPromotions() (model.Pagination, error)
}

type promotionService struct {
	repo repository.PromotionRepository
}

func NewPromotionService(
	repo repository.PromotionRepository,
) PromotionService {
	return &promotionService{repo}
}

func (s *promotionService) GetPromotions() (model.Pagination, error) {

	promotions, err := s.repo.GetPromotions()
	if err != nil {
		return model.Pagination{}, err
	}

	return model.Pagination{
		List:  promotions,
		Total: 0,
	}, nil
}
