package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &repo{db}
}

type GroupRepository interface {
	Create(data *model.CreateGroup) error
}

func (r repo) Create(data *model.CreateGroup) error {

	if err := r.db.Table("Groups").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}
