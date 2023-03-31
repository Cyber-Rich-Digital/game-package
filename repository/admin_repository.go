package repository

import (
	"cybergame-api/model"
	"time"

	"gorm.io/gorm"
)

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &repo{db}
}

type AdminRepository interface {
	GetGroup(groupId int) (*model.AdminGroupPermissionResponse, error)
	GetGroupList(query model.AdminGroupQuery) (*[]model.GroupList, *int64, error)
	GetAdminByUsername(data model.LoginAdmin) (*model.Admin, error)
	CheckUsername(username string) (bool, error)
	CheckPhone(phone string) (bool, error)
	CreateAdmin(user model.Admin) error
	CreateGroup(data []model.AdminPermissionList) error
	UpdateGroup(data []model.AdminPermissionList, perIds []int64) error
}

func (r repo) GetGroup(groupId int) (*model.AdminGroupPermissionResponse, error) {

	var group model.Group
	var permission []model.PermissionList

	if err := r.db.Table("Admin_groups").
		Select("id, name").
		Where("id = ?", groupId).
		First(&group).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.Table("Permissions p").
		Select("p.id, p.name").
		Joins("LEFT JOIN Admin_group_permissions gp ON gp.permission_id = p.id").
		Where("gp.group_id = ?", groupId).
		Find(&permission).
		Error; err != nil {
		return nil, err
	}

	var result model.AdminGroupPermissionResponse
	result.Id = group.Id
	result.Name = group.Name
	result.Permissions = permission

	return &result, nil
}

func (r repo) GetGroupList(query model.AdminGroupQuery) (*[]model.GroupList, *int64, error) {

	var list []model.GroupList
	if err := r.db.Table("Admin_groups").
		Select("id, name, admin_count").
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&list).
		Error; err != nil {
		return nil, nil, err
	}

	var total int64
	if err := r.db.Table("Admin_groups").
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &list, &total, nil
}

func (r repo) GetAdminByUsername(data model.LoginAdmin) (*model.Admin, error) {
	var admin model.Admin

	if err := r.db.Table("Admins").
		Select("id, username, phone, password, email, role").
		Where("username = ?", data.Username).
		First(&admin).
		Error; err != nil {
		return nil, err
	}

	if admin.Id != 0 {
		if err := r.db.Table("Admins").
			Where("id = ?", admin.Id).
			Updates(model.AdminLoginUpdate{
				IP:        data.IP,
				LogedinAt: time.Now(),
			}).
			Error; err != nil {
			return nil, err
		}
	}

	return &admin, nil
}

func (r repo) CheckUsername(username string) (bool, error) {
	var user model.Admin

	if err := r.db.Table("Admins").
		Where("username = ?", username).
		First(&user).
		Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r repo) CheckPhone(phone string) (bool, error) {
	var user model.Admin

	if err := r.db.Table("Admins").
		Where("phone = ?", phone).
		First(&user).
		Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r repo) CreateAdmin(user model.Admin) error {

	if err := r.db.Table("Admins").
		Create(&user).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) CreateGroup(data []model.AdminPermissionList) error {

	if err := r.db.Table("Admin_group_permissions").
		Create(&data).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateGroup(data []model.AdminPermissionList, perIds []int64) error {

	tx := r.db.Begin()

	if err := tx.Table("Admin_group_permissions").
		Where("group_id = ? AND permission_id IN (?)", data[0].GroupId, perIds).
		Delete(&model.AdminGroupPermission{}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table("Admin_group_permissions").
		Create(&data).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
