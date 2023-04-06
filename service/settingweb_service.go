package service

import (
	"cybergame-api/model"
	"cybergame-api/repository"

	"cloud.google.com/go/storage"
)

var (
	storageClient *storage.Client
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

	var web model.Settingweb
	web.Logo = data.Logo
	web.BackgrondColor = data.UserAuto
	web.UserAuto = data.UserAuto
	web.OtpRegister = data.OtpRegister
	web.TranWithdraw = data.TranWithdraw
	web.Register = data.Register
	web.DepositFirst = data.DepositFirst
	web.DepositNext = data.DepositNext
	web.Withdraw = data.Withdraw
	web.Line = data.Line
	web.Url = data.Url
	web.Opt = data.Opt

	if err := s.repo.CreateSettingWeb(data); err != nil {
		return internalServerError(err.Error())
	}
	return nil
}
