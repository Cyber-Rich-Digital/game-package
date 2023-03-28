package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

type DeviceService interface {
	CreateDevice(data *model.DeviceBody) (*string, error)
}

type deviceService struct {
	repo        repository.DeviceRepository
	websiteRepo repository.WebsiteRepository
}

func NewDeviceService(
	repo repository.DeviceRepository,
	websiteRepo repository.WebsiteRepository,
) DeviceService {
	return &deviceService{repo, websiteRepo}
}

func (s *deviceService) CreateDevice(data *model.DeviceBody) (*string, error) {

	exist, err := s.repo.CheckDevice(*data)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound("Device not found")
		}

		return nil, internalServerError(err.Error())
	}

	website, err := s.websiteRepo.GetWebsite(data.WebsiteId)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound("Device not found")
		}

		return nil, internalServerError(err.Error())
	}

	var deviceId *int

	if exist == nil {

		id, err := s.repo.CreateDevice(*data)
		if err != nil {
			return nil, internalServerError(err.Error())
		}

		deviceId = id
	} else {

		deviceId = exist

		if err := s.repo.UpdateDevice(*deviceId, *data); err != nil {
			return nil, internalServerError(err.Error())
		}
	}

	token, err := helper.CreateJWT(data.HardwareId, *deviceId, website.UserId)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}
