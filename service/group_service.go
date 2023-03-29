package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
)

type GroupService interface {
	Create(user *model.CreateGroup) error
}

const GroupNotFound = "Permission not found"

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(
	repo repository.GroupRepository,
) GroupService {
	return &groupService{repo}
}

func (s *groupService) Create(data *model.CreateGroup) error {

	if err := s.repo.Create(data); err != nil {
		return err
	}

	return nil
}
