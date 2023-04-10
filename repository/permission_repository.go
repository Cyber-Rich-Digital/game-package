package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &repo{db}
}

type PermissionRepository interface {
	CheckPerListExist(ids []int64) ([]int64, error)
	CheckPerListAndGroupId(groupId int64, perIds []int64) ([]int64, error)
	CreatePermission(data *model.CreatePermission) error
	DeletePermission(perm model.DeletePermission) error
}

func (r repo) CheckPerListExist(ids []int64) ([]int64, error) {

	var list []int64

	if err := r.db.Table("Permissions").Select("id").Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) CheckPerListAndGroupId(groupId int64, perIds []int64) ([]int64, error) {

	var list []int64

	if err := r.db.Table("Admin_group_permissions").
		Select("id").
		Where("group_id = ? AND permission_id IN (?)", groupId, perIds).
		Find(&list).
		Error; err != nil {
		return nil, err
	}

	return list, nil
}

func (r repo) CreatePermission(data *model.CreatePermission) error {

	if err := r.db.Table("Permissions").
		Create(&data.Permissions).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeletePermission(perm model.DeletePermission) error {

	tx := r.db.Begin()

	if err := tx.Table("Permissions").
		Where("id IN (?)", perm.PermissionIds).
		Delete(&model.Permission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Where("permission_id IN (?)", perm.PermissionIds).
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
