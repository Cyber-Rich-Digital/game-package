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
	DeleteGroup(id int64) error
}

func (r repo) CheckGroupExist(id int64) (bool, error) {

	var count int64

	if err := r.db.Table("Admin_groups").
		Where("id = ?", id).
		Count(&count).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

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

func (r repo) DeleteGroup(id int64) error {

	tx := r.db.Begin()

	if err := tx.Table("Admin_groups").
		Where("id = ?", id).
		Delete(&model.Group{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Where("group_id = ?", id).
		Delete(&model.AdminGroupPermission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
