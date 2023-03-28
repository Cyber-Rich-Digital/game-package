package repository

import (
	"cybergame-api/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repo{db}
}

type UserRepository interface {
	GetUsers(query model.UserQuery) (*[]model.UserResponse, int64, error)
	GetAdmins(data model.UserQuery) (*model.Pagination, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	CheckUserByEmailOrUser(val string) (bool, error)
	CheckRole() (bool, error)
	CreateUser(user model.User) error
	CreateAdmin(user model.User) error
	ChangePassword(id int, password string) error
	DeleteUser(id int) error
}

func (r repo) GetUsers(data model.UserQuery) (*[]model.UserResponse, int64, error) {

	var users *[]model.UserResponse
	var total int64
	var err error

	selectFields := "u.id, u.username, u.email, u.created_at, COUNT(w.id) AS web_total"
	join := "LEFT JOIN Websites AS w ON u.id = w.user_id"
	group := "w.user_id, u.id, u.username, u.email, u.created_at"
	whereVal := fmt.Sprintf("%%%s%%", data.Search)

	// Get list of users //

	query := r.db.Table("Users u")

	query = query.
		Select(selectFields).
		Joins(join).
		Group(group)

	if data.Search != "" {
		query = query.Where("u.email LIKE ?", whereVal).Or("u.username LIKE ?", whereVal)
	}

	if data.Sort == 1 {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err = query.
		Where("u.role = ?", "USER").
		Where("u.deleted_at IS NULL").
		Limit(data.Limit).
		Offset(data.Page * data.Limit).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// Get total count //

	count := r.db.Table("Users").
		Select("id, email, username")

	if data.Search != "" {
		count = count.Where("email LIKE ?", whereVal).Or("username LIKE ?", whereVal)
	}

	if err = count.
		Where("role = ?", "USER").
		Where("deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, 0, err
	}

	// Return response //

	return users, total, nil
}

func (r repo) GetAdmins(data model.UserQuery) (*model.Pagination, error) {

	var users *[]model.UserAdminResponse
	var total int64
	var err error

	selectFields := "id, username, email, created_at"
	whereVal := fmt.Sprintf("%%%s%%", data.Search)

	// Get list of users //

	query := r.db.Table("Users")

	query = query.
		Select(selectFields)

	if data.Search != "" {
		query = query.Where("email LIKE ?", whereVal).Or("username LIKE ?", whereVal)
	}

	if data.Sort == 1 {
		query = query.Order("created_at ASC")
	} else {
		query = query.Order("created_at DESC")
	}

	if err = query.
		Where("role IN (?)", []string{"ADMIN", "SUPER_ADMIN"}).
		Where("deleted_at IS NULL").
		Limit(data.Limit).
		Offset(data.Page * data.Limit).
		Find(&users).Error; err != nil {
		return nil, err
	}

	// Get total count //

	count := r.db.Table("Users").
		Select("id")

	if data.Search != "" {
		count = count.Where("email LIKE ?", whereVal).Or("username LIKE ?", whereVal)
	}

	if err = count.
		Where("role = ?", "ADMIN").
		Where("deleted_at IS NULL").
		Count(&total).
		Error; err != nil {
		return nil, err
	}

	// Return response //

	return &model.Pagination{
		List:  users,
		Total: total,
	}, nil

}

func (r repo) GetUserByID(id int) (*model.User, error) {

	var user model.User
	if err := r.db.Table("Users").
		Select("id, email, password, username").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First(&user).
		Error; err != nil {
		return nil, err
	}

	if user.Id == 0 {
		return nil, nil
	}

	return &user, nil
}

func (r repo) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Table("Users").
		Select("id, email, username, password, role").
		Where("email = ?", email).
		Or("username = ?", email).
		Where("deleted_at IS NULL").
		First(&user).
		Error; err != nil {
		return nil, err
	}

	if user.Id == 0 {
		return nil, errors.New("User not found")
	}

	return &user, nil
}

func (r repo) CheckUserByEmailOrUser(val string) (bool, error) {

	var result int64
	if err := r.db.Table("Users").
		Select("id").
		Where("username = ?", val).
		Or("email = ?", val).
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&result).
		Error; err != nil {
		return false, err
	}

	return result > 0, nil
}

func (r repo) CheckRole() (bool, error) {

	var result int64
	if err := r.db.Table("Users").
		Select("id").
		Where("role = ?", "ADMIN").
		Where("deleted_at IS NULL").
		Limit(1).
		Count(&result).
		Error; err != nil {
		return false, err
	}

	return result > 0, nil
}

func (r repo) CreateUser(user model.User) error {
	if err := r.db.Table("Users").
		Create(&user).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) CreateAdmin(user model.User) error {
	if err := r.db.Table("Users").
		Create(&user).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) ChangePassword(id int, password string) error {
	if err := r.db.Table("Users").
		Where("id = ?", id).
		Update("password", password).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteUser(id int) error {

	if err := r.db.Table("Users").
		Where("id = ?", id).
		Delete(&model.User{}).
		Error; err != nil {
		return err
	}

	return nil
}
