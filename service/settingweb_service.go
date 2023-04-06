package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
)

func NewSettingWebService(
	repo repository.SettingWebRepository,
) SettingWebService {
	return &settingwebService{repo}
}

type SettingWebService interface {
	CreateSettingWeb(data model.SettingwebCreateBody) error
}

type settingwebService struct {
	repo repository.SettingWebRepository
}

// CreateSettingWeb implements SettingWebService
func (s *settingwebService) CreateSettingWeb(data model.SettingwebCreateBody) error {
	if err := s.repo.CreateSettingWeb(data); err != nil {
		return err
	}

	return nil
}
