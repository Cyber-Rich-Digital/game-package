package service

import (
	"cyber-api/helper"
	"cyber-api/model"
	"cyber-api/repository"

	"golang.org/x/crypto/bcrypt"
)

var loginFailed = "Email Or Password is incorrect"

type UserService interface {
	GetUserByID(id int) (*model.UserReponse, error)
	CreateUser(user *model.CreateUser) error
	Login(user *model.Login) (*string, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(
	repo repository.UserRepository,
) UserService {
	return &userService{repo}
}

func (s *userService) Login(user *model.Login) (*string, error) {

	exist, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	if exist.Id == 0 {
		return nil, notFound("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.Password)); err != nil {
		return nil, notFound(loginFailed)
	}

	token, err := helper.CreateJWT(&exist)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *userService) GetUserByID(id int) (*model.UserReponse, error) {

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return user, nil
}

func (s *userService) CreateUser(user *model.CreateUser) error {

	exist, err := s.repo.CheckUserByEmail(user.Email)
	if err != nil {
		return err
	}

	if exist {
		return badRequest("Email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError(err.Error())
	}

	user.Password = string(hashedPassword)

	return s.repo.CreateUser(user)
}
