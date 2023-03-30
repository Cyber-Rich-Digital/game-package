package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &repo{db}
}

type GroupRepository interface {
	CheckGroupExist(id int64) (bool, error)
	Create(data *model.CreateGroup) error
}

func (r repo) CheckGroupExist(id int64) (bool, error) {

	var count int64

	if err := r.db.Table("Admin_groups").Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) Create(data *model.CreateGroup) error {

	if err := r.db.Table("Admin_groups").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}
