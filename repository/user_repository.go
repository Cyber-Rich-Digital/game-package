package repository

import (
	"cyber-api/model"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repo{db}
}

type UserRepository interface {
	GetUserByID(id int) (model.User, error)
	GetUserByEmail(email string) (model.User, error)
	GetUsers() (interface{}, error)
	CheckUserByEmail(email string) (bool, error)
	CreateUser(user *model.CreateUser) error
}

func (r repo) GetUserByID(id int) (model.User, error) {

	var user model.User
	if err := r.db.Table("Users").
		Select("id, email, password, username").
		Where("id = ?", id).
		First(&user).
		Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r repo) GetUserByEmail(email string) (model.User, error) {

	var user model.User
	if err := r.db.Table("Users").
		Select("id, email, password, username").
		Where("email = ?", email).
		First(&user).
		Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r repo) GetUsers() (interface{}, error) {

	var promotions interface{}
	if err := r.db.Table("Users").Find(&promotions).Error; err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r repo) CheckUserByEmail(email string) (bool, error) {

	var result int64
	if err := r.db.Table("Users").Where("email = ?", email).Count(&result).Error; err != nil {
		return false, err
	}

	return result > 0, nil
}

func (r repo) CreateUser(user *model.CreateUser) error {
	if err := r.db.Table("Users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}
