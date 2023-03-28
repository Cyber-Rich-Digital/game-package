package service

import "cybergame-api/model"

type MenuService interface {
	GetMenu(role string) []model.Menu
}

type menuService struct {
}

func NewMenuService() MenuService {
	return &menuService{}
}

func (s *menuService) GetMenu(role string) []model.Menu {

	var list []model.Menu

	if role != "USER" {
		list = []model.Menu{
			{Menu: "Websites", Path: "/"},
			{Menu: "Users", Path: "/users"},
			{Menu: "Admins", Path: "/admins"},
		}
	} else {
		list = []model.Menu{
			{Menu: "Websites", Path: "/websites"},
		}
	}

	return list
}
