package repository

import (
	"cybergame-api/model"

	"gorm.io/gorm"
)

func NewUserRepository(db *gorm.DB) UserRepository {
	return &repo{db}
}

type UserRepository interface {
	GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error)
	GetUser(id int64) (*model.UserDetail, error)
	GetUserList(query model.UserListQuery) (*[]model.UserList, *int64, error)
	CheckUser(username string) (bool, error)
	CheckUserPhone(phone string) (bool, error)
	CheckUserById(id int64) (bool, error)
	GetUserByPhone(phone string) (*model.UserByPhone, error)
	CreateUser(admin model.User) error
	UpdateUser(userId int64, data model.UpdateUser) error
	UpdateUserPassword(userId int64, data model.UserUpdatePassword) error
	DeleteUser(id int64) error
}

func (r repo) GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error) {

	var logs []model.UserLoginLog

	if err := r.db.Table("User_login_logs").
		Where("user_id = ?", id).
		Find(&logs).
		Order("created_at DESC").
		Error; err != nil {
		return nil, err
	}

	return &logs, nil
}

func (r repo) GetUserList(query model.UserListQuery) (*[]model.UserList, *int64, error) {

	var err error
	var list []model.UserList
	var total int64

	exec := r.db.Table("Users").
		Select("id, member_code, promotion, fullname, bankname, bank_account, channel, credit, ip, ip_registered, created_at, updated_at, logedin_at")

	if query.Search != "" {
		exec = exec.Where("username LIKE ?", "%"+query.Search+"%")
	}

	if query.Status != "" {
		exec = exec.Where("status = ?", query.Status)
	}

	if err := exec.
		Limit(query.Limit).
		Offset(query.Limit * query.Page).
		Find(&list).
		Error; err != nil {
		return nil, nil, err
	}

	execTotal := r.db.Table("Users").
		Select("id")

	if query.Search != "" {
		execTotal = execTotal.Where("username LIKE ?", "%"+query.Search+"%")
	}

	if query.Status != "" {
		execTotal = execTotal.Where("status = ?", query.Status)
	}

	if err = execTotal.
		Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	return &list, &total, nil
}

func (r repo) GetUser(id int64) (*model.UserDetail, error) {

	var admin *model.UserDetail

	if err := r.db.Table("Users").
		Select("id, partner, member_code, phone, promotion, bankname, bank_account, fullname, channel, true_wallet, contact, note, course").
		Where("id = ?", id).
		First(&admin).
		Error; err != nil {
		return nil, err
	}

	return admin, nil
}

func (r repo) CheckUser(username string) (bool, error) {
	var user model.User

	if err := r.db.Table("Users").
		Where("username = ?", username).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if user.Id != 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CheckUserPhone(phone string) (bool, error) {

	var user model.User

	if err := r.db.Table("Users").
		Where("phone = ?", phone).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	if user.Id != 0 {
		return false, nil
	}

	return true, nil
}

func (r repo) CheckUserById(id int64) (bool, error) {
	var user model.User

	if err := r.db.Table("Users").
		Where("id = ?", id).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (r repo) GetUserByPhone(phone string) (*model.UserByPhone, error) {

	var user *model.UserByPhone

	if err := r.db.Table("Users").
		Select("id, phone").
		Where("phone = ?", phone).
		First(&user).
		Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return user, nil
}

func (r repo) CreateUser(admin model.User) error {

	if err := r.db.Table("Users").
		Create(&admin).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) UpdateUser(userId int64, data model.UpdateUser) error {

	if err := r.db.Table("Users").
		Where("id = ?", userId).
		Updates(&data).
		Error; err != nil {
		r.db.Rollback()
		return err
	}

	return nil
}

func (r repo) UpdateUserPassword(userId int64, data model.UserUpdatePassword) error {

	if err := r.db.Table("Users").
		Where("id = ?", userId).
		Update("password", data.Password).
		Error; err != nil {
		return err
	}

	return nil
}

func (r repo) DeleteUser(id int64) error {

	if err := r.db.Table("Users").
		Where("id = ?", id).
		Delete(&model.User{}).
		Error; err != nil {
		return err
	}

	return nil
}
