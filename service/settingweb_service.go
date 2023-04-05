package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"
)

type SettingWebService interface {
	GetSettingWeb(data model.SettingwebListRequest) (*model.SuccessWithPagination, error)
	GetSettingWebById(data model.SettingwebParam) (*model.Settingweb, error)
	CreateSettingWeb(data model.SettingwebCreateBody) error
	UpdateSettingWeb(id int64, data model.SettingwebUpdateBody) error
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

// GetSettingWeb implements SettingWebService
func (*settingwebService) GetSettingWeb(data model.SettingwebListRequest) (*model.SuccessWithPagination, error) {
	panic("unimplemented")
}

// GetSettingWebById implements SettingWebService
func (*settingwebService) GetSettingWebById(data model.SettingwebParam) (*model.Settingweb, error) {
	panic("unimplemented")
}

// UpdateSettingWeb implements SettingWebService
func (*settingwebService) UpdateSettingWeb(id int64, data model.SettingwebUpdateBody) error {
	panic("unimplemented")
}

func NewSettingWebService(
	repo repository.SettingWebRepository,
) SettingWebService {
	return &settingwebService{repo}
}
