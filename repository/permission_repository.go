package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &repo{db}
}

type PermissionRepository interface {
	CreatePermission(data *model.CreatePermission) error
}

func (r repo) CreatePermission(data *model.CreatePermission) error {

	if err := r.db.Table("Permissions").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}
