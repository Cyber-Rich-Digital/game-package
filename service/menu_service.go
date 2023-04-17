package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
)

type MenuService interface {
	GetMenu() ([]model.Menu, error)
}

type menuService struct {
	PermRepo repository.PermissionRepository
}

func NewMenuService(
	PermRepo repository.PermissionRepository,
) MenuService {
	return &menuService{PermRepo}
}

func (s *menuService) GetMenu() ([]model.Menu, error) {

	perms, err := s.PermRepo.GetPermissions()
	if err != nil {
		return nil, err
	}

	var menu []model.Menu

	for _, perm := range perms {
		menu = append(menu, model.Menu{
			Id:    perm.Id,
			Title: perm.Name,
			Name:  perm.Name,
			View:  false,
			Edit:  false,
			List:  nil,
		})
	}

	return menu, nil
}
