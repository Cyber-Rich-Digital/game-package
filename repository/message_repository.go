package repository

import (
	"cyber-api/model"

	"gorm.io/gorm"
)

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &repo{db}
}

type MessageRepository interface {
	GetMessagesByTagId(id int) (*[]model.MessageResponse, error)
	CreateMessage(data model.MessageCreate) error
	CreateMessageAndTag(data model.MessageCreate, tagData model.TagInsert) error
}

func (r repo) GetMessagesByTagId(id int) (*[]model.MessageResponse, error) {

	var messages []model.MessageResponse

	if err := r.db.Table("Messages").
		Select("id, message, created_at").
		Where("tag_id = ?", id).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	return &messages, nil
}

func (r repo) CreateMessage(data model.MessageCreate) error {

	if err := r.db.Table("Messages").Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r repo) CreateMessageAndTag(data model.MessageCreate, tagData model.TagInsert) error {

	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return tx.Error
	}

	if err := tx.Table("Messages").Create(&data).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Tags").Create(&tagData).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
