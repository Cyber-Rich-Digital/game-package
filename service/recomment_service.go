package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

type RecommentService interface {
	GetRecommentList(query model.RecommentQuery) ([]model.RecommentList, int64, error)
	CreateRecomment(user model.CreateRecomment) error
	UpdateRecomment(id int64, body model.CreateRecomment) error
}

type recommentService struct {
	repo repository.RecommentRepository
}

func NewRecommentService(
	repo repository.RecommentRepository,
) RecommentService {
	return &recommentService{repo}
}

func (s *recommentService) GetRecommentList(query model.RecommentQuery) ([]model.RecommentList, int64, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, 0, err
	}

	return s.repo.GetRecommentList(query)
}

func (s *recommentService) CreateRecomment(body model.CreateRecomment) error {

	if err := s.repo.CreateRecomment(body); err != nil {
		return err
	}

	return nil
}

func (s *recommentService) UpdateRecomment(id int64, body model.CreateRecomment) error {

	if err := s.repo.UpdateRecomment(id, body); err != nil {
		return err
	}

	return nil
}
