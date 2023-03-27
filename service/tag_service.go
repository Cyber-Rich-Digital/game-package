package service

import (
	"cyber-api/model"
	"cyber-api/repository"
)

type TagService interface {
	GetTagsByWebsiteId(data model.TagParam) (*[]model.TagResponse, error)
	DeleteTag(id int) error
}

type tagService struct {
	repo repository.TagRepository
}

func NewTagService(
	repo repository.TagRepository,
) TagService {
	return &tagService{repo}
}

func (s *tagService) GetTagsByWebsiteId(data model.TagParam) (*[]model.TagResponse, error) {

	tags, err := s.repo.GetTagsByWebsiteId(data)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	if tags == nil {
		return nil, notFound("Tags not found")
	}

	return tags, nil
}

func (s *tagService) DeleteTag(id int) error {

	if err := s.repo.DeleteTag(id); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
