package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
)

type AdminService interface {
	Login(data model.LoginAdmin) (*string, error)
	Create(user *model.CreateAdmin) error
}

var AdminloginFailed = "Phone Or Password is incorrect"

const AdminNotFound = "Admin not found"

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(
	repo repository.AdminRepository,
) AdminService {
	return &adminService{repo}
}

func (s *adminService) Login(data model.LoginAdmin) (*string, error) {

	user, err := s.repo.GetAdminByUsername(data)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(AdminloginFailed)
		}

		return nil, internalServerError(err.Error())
	}

	if user == nil {
		return nil, badRequest(AdminloginFailed)
	}

	if err := helper.ComparePassword(data.Password, user.Password); err != nil {
		return nil, badRequest(AdminloginFailed)
	}

	token, err := helper.CreateJWT(*user)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *adminService) Create(data *model.CreateAdmin) error {

	username, err := s.repo.CheckUsername(data.Username)
	if err != nil {
		return err
	}

	if username {
		return badRequest("Username already exist")
	}

	email, err := s.repo.CheckPhone(data.Phone)
	if err != nil {
		return err
	}

	if email {
		return badRequest("Phone already exist")
	}

	hashedPassword, err := helper.GenPassword(data.Password)
	if err != nil {
		return internalServerError(err.Error())
	}

	newUser := model.Admin{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Password = string(hashedPassword)
	newUser.Role = "ADMIN"

	return s.repo.CreateAdmin(newUser)
}
