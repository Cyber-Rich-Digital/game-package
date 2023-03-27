package repository

import (
	"cyber-api/model"

	"gorm.io/gorm"
)

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &repo{db}
}

type DeviceRepository interface {
	CheckDevice(data model.DeviceBody) (*int, error)
	GetTokensbyWebsiteId(websiteId int) (*[]model.Device, error)
	GetDevicesByWebId(webId int) (*[]model.Device, error)
	CreateDevice(data model.DeviceBody) (*int, error)
	UpdateDevice(id int, data model.DeviceBody) error
}

func (r repo) CheckDevice(data model.DeviceBody) (*int, error) {

	var device model.Device

	if err := r.db.Table("Devices").
		Select("id").
		Where("hardware_id = ?", data.HardwareId).
		Limit(1).
		Find(&device).
		Error; err != nil {
		return nil, err
	}

	if device.Id == 0 {
		return nil, nil
	}

	return &device.Id, nil
}

func (r repo) GetDevicesByWebId(webId int) (*[]model.Device, error) {

	var list *[]model.Device
	if err := r.db.Table("Devices").
		Select("id").
		Where("website_id = ?", webId).
		Find(&list).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) GetTokensbyWebsiteId(websiteId int) (*[]model.Device, error) {

	var devices []model.Device
	if err := r.db.Table("Devices").
		Select("fcm_token").
		Where("website_id = ?", websiteId).
		Find(&devices).
		Error; err != nil {
		return nil, err
	}

	return &devices, nil
}

func (r repo) CreateDevice(data model.DeviceBody) (*int, error) {

	device := model.Device{
		WebsiteId:  data.WebsiteId,
		HardwareId: data.HardwareId,
		FcmToken:   data.FcmToken,
	}

	result := r.db.Table("Devices").Create(&device)
	if result.Error != nil {
		return nil, result.Error
	}

	return &device.Id, nil
}

func (r repo) UpdateDevice(id int, data model.DeviceBody) error {

	result := r.db.Table("Devices").Where("id = ?", id).Updates(map[string]interface{}{
		"fcm_token":  data.FcmToken,
		"website_id": data.WebsiteId,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
