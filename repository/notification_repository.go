package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewNotiRepository(db *gorm.DB) NotiRepository {
	return &repo{db}
}

type NotiRepository interface {
	GetNotiByTagIdAndDeviceId(tagId int, deviceId int) (*model.Notifications, error)
	CreateNoti(data model.Notifications) error
	UpdateNoti(tagId int, deviceId int) error
	SetRead(data model.MessageRead) error
	// CreateMessageAndTag(data model.MessageCreate, tagData model.TagInsert) error
}

func (r repo) GetNotiByTagIdAndDeviceId(tagId int, deviceId int) (*model.Notifications, error) {

	var messages model.Notifications

	if err := r.db.Table("Notifications").
		Select("id").
		Where("tag_id = ?", tagId).
		Where("device_id = ?", deviceId).
		First(&messages).
		Error; err != nil {
		return nil, err
	}

	return &messages, nil
}

func (r repo) CreateNoti(data model.Notifications) error {

	if err := r.db.Table("Notifications").Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateNoti(tagId int, deviceId int) error {

	if err := r.db.Table("Notifications").
		Where("tag_id = ?", tagId).
		Where("device_id = ?", deviceId).
		Update("total", gorm.Expr("total + ?", 1)).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) SetRead(data model.MessageRead) error {

	if err := r.db.Table("Notifications").
		Where("tag_id = ?", data.TagId).
		Where("device_id = ?", data.DeviceId).
		Update("total", 0).
		Error; err != nil {
		return err
	}

	return nil
}

// func (r repo) CreateMessageAndTag(data model.MessageCreate, tagData model.TagInsert) error {

// 	tx := r.db.Begin()

// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	if err := tx.Error; err != nil {
// 		return tx.Error
// 	}

// 	if err := tx.Table("Messages").Create(&data).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Table("Tags").Create(&tagData).Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		tx.Rollback()
// 		return err
// 	}

// 	return nil
// }
