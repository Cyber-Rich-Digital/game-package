package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"
	"reflect"
	"strings"
)

type UserService interface {
	GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error)
	GetUser(id int64) (*model.UserDetail, error)
	GetUserList(query model.UserListQuery) (*model.SuccessWithPagination, error)
	Create(user *model.CreateUser) error
	UpdateUser(userId int64, body model.UpdateUser, adminName string) error
	ResetPassword(userId int64, body model.UserUpdatePassword) error
	DeleteUser(id int64) error
}

const UserloginFailed = "Phone Or Password is incorrect"
const UserNotFound = "User not found"
const UserExist = "User already exist"
const UserPhoneExist = "Phone already exist"

type userService struct {
	repo      repository.UserRepository
	perRepo   repository.PermissionRepository
	groupRepo repository.GroupRepository
}

func NewUserService(
	repo repository.UserRepository,
	perRepo repository.PermissionRepository,
	groupRepo repository.GroupRepository,
) UserService {
	return &userService{repo, perRepo, groupRepo}
}

func (s *userService) GetUserLoginLogs(id int64) (*[]model.UserLoginLog, error) {

	logs, err := s.repo.GetUserLoginLogs(id)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *userService) GetUser(id int64) (*model.UserDetail, error) {

	admin, err := s.repo.GetUser(id)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(UserNotFound)
		}

		return nil, err
	}

	return admin, nil
}

func (s *userService) GetUserList(query model.UserListQuery) (*model.SuccessWithPagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, err
	}

	list, total, err := s.repo.GetUserList(query)
	if err != nil {
		return nil, err
	}

	result := &model.SuccessWithPagination{
		Message: "Success",
		List:    list,
		Total:   *total,
	}

	return result, nil
}

func (s *userService) Create(data *model.CreateUser) error {

	userByPhone, err := s.repo.CheckUserPhone(data.Phone)
	if err != nil {
		return err
	}

	if userByPhone {
		return badRequest("Phone already exist")
	}

	hashedPassword, err := helper.GenUserPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	var newUser model.User
	newUser.Partner = data.Partner
	newUser.MemberCode = data.MemberCode
	newUser.Username = data.Phone
	newUser.Phone = data.Phone
	newUser.Promotion = data.Promotion
	newUser.Password = string(hashedPassword)
	newUser.Status = data.Status
	newUser.Fullname = data.Fullname
	newUser.Bankname = data.Bankname
	newUser.BankAccount = data.BankAccount
	newUser.Channel = data.Channel
	newUser.TrueWallet = data.TrueWallet
	newUser.Contact = data.Contact
	newUser.Note = data.Note
	newUser.Course = data.Course
	newUser.IpRegistered = data.IpRegistered

	splitFullname := strings.Split(data.Fullname, " ")
	var firstname, lastname *string
	if len(splitFullname) > 1 {
		firstname = &splitFullname[0]
		lastname = &splitFullname[1]
		newUser.Firstname = *firstname
		newUser.Lastname = *lastname
	}

	return s.repo.CreateUser(newUser)
}

func (s *userService) UpdateUser(userId int64, body model.UpdateUser, adminName string) error {

	user, err := s.repo.GetUser(userId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if user == nil {
		return notFound(UserNotFound)
	}

	var changeList []model.UserUpdateLogs

	b := reflect.ValueOf(body)
	u := reflect.ValueOf(user)

	if b.Kind() == reflect.Ptr {
		b = b.Elem()
	}

	if u.Kind() == reflect.Ptr {
		u = u.Elem()
	}

	// loop user fields
	for j := 0; j < b.NumField(); j++ {
		for k := 0; k < u.NumField(); k++ {

			bField := b.Type().Field(j).Name
			bValue := b.Field(j).Interface()
			uField := u.Type().Field(k).Name
			uValue := u.Field(k).Interface()

			if bField == uField {
				if uValue != bValue {
					changeList = append(changeList, model.UserUpdateLogs{
						UserId:            userId,
						Description:       fmt.Sprintf("%s changed from %s to %s", bField, uValue, bValue),
						CreatedByUsername: adminName,
						Ip:                body.Ip,
					})
				}
			}
		}
	}

	return s.repo.UpdateUser(userId, body, changeList)
}

func (s *userService) ResetPassword(userId int64, body model.UserUpdatePassword) error {

	checkUser, err := s.repo.CheckUserById(userId)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkUser {
		return notFound(UserNotFound)
	}

	newPasword, err := helper.GenUserPassword(body.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	body.Password = newPasword

	if err := s.repo.UpdateUserPassword(userId, body); err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(id int64) error {

	checkUser, err := s.repo.CheckUserById(id)
	if err != nil {
		return internalServerError(err.Error())
	}

	if !checkUser {
		return notFound(UserNotFound)
	}

	return s.repo.DeleteUser(id)
}
