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

	var group model.Group
	group.Name = data.Name

	var perms []model.Permission

	for _, permission := range data.Permissions {
		perms = append(perms, model.Permission{
			Name: permission.Name,
		})
	}

	if err := s.repo.Create(group, perms); err != nil {
		return err
	}

	return nil
}
