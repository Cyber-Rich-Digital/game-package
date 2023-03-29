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
	GetAdminByUsername(data model.LoginAdmin) (*model.Admin, error)
	CheckUsername(username string) (bool, error)
	CheckPhone(phone string) (bool, error)
	CreateAdmin(user model.Admin) error
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
