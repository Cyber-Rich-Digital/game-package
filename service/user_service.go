package service

import (
	"cybergame-api/helper"
	"cybergame-api/model"
	"cybergame-api/repository"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var loginFailed = "Username Or Password is incorrect"

type UserService interface {
	GetUsers(query model.UserQuery) (*model.Pagination, error)
	GetAdmins(query model.UserQuery) (*model.Pagination, error)
	Login(user *model.Login) (*string, error)
	CreateUser(user *model.CreateUser) error
	CreateAdmin(user *model.CreateAdmin, admin bool) error
	UserChangePassword(userId1, userId2 int, role string, user *model.UserChangePassword) error
	AdminChangePassword(id int, user *model.AdminChangePassword) error
	DeleteUser(id int) error
}

type userService struct {
	repo        repository.UserRepository
	websiteRepo repository.WebsiteRepository
}

func NewUserService(
	repo repository.UserRepository,
	websiteRepo repository.WebsiteRepository,
) UserService {
	return &userService{repo, websiteRepo}
}

func (s *userService) GetUsers(query model.UserQuery) (*model.Pagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	users, total, err := s.repo.GetUsers(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	userIds := make([]int, 0)
	for _, user := range *users {
		userIds = append(userIds, int(user.Id))
	}

	websites, err := s.websiteRepo.GetWebsitesByUserIds(userIds)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	for i, user := range *users {
		for _, website := range *websites {
			if user.Id == website.UserId {
				(*users)[i].Websites = append((*users)[i].Websites, website)
			}
		}
	}

	returnUsers := &model.Pagination{
		Total: total,
		List:  users,
	}

	return returnUsers, nil
}

func (s *userService) GetAdmins(query model.UserQuery) (*model.Pagination, error) {

	if err := helper.Pagination(&query.Page, &query.Limit); err != nil {
		return nil, badRequest(err.Error())
	}

	users, err := s.repo.GetAdmins(query)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return users, nil
}

func (s *userService) Login(user *model.Login) (*string, error) {

	exist, err := s.repo.GetUserByEmail(user.Email)
	if err != nil {

		if err.Error() == "record not found" {
			return nil, notFound(loginFailed)
		}

		if err.Error() == "User not found" {
			return nil, notFound(loginFailed)
		}

		return nil, internalServerError(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.Password)); err != nil {
		return nil, notFound(loginFailed)
	}

	token, err := helper.CreateJWTUser(exist)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return &token, nil
}

func (s *userService) GetUserByID(id int) (*model.User, error) {

	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, internalServerError(err.Error())
	}

	return user, nil
}

func (s *userService) CreateUser(data *model.CreateUser) error {

	email, emailErr := s.repo.CheckUserByEmailOrUser(data.Email)
	if emailErr != nil {
		return emailErr
	}

	if email {
		return badRequest("Email already exist")
	}

	user, userErr := s.repo.CheckUserByEmailOrUser(data.Username)
	if userErr != nil {
		return userErr
	}

	if user {
		return badRequest("Username already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError(err.Error())
	}

	newUser := model.User{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Password = string(hashedPassword)
	newUser.Role = "USER"

	return s.repo.CreateUser(newUser)
}

func (s *userService) CreateAdmin(data *model.CreateAdmin, admin bool) error {

	if !admin {

		role, err := s.repo.CheckRole()

		if err != nil {
			return err
		}

		if role {
			return badRequest("Admin already exist")
		}

	}

	user, err := s.repo.CheckUserByEmailOrUser(data.Username)
	if err != nil {
		return err
	}

	if user {
		return badRequest("Username already exist")
	}

	email, err := s.repo.CheckUserByEmailOrUser(data.Email)
	if err != nil {
		return err
	}

	if email {
		return badRequest("Email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError(err.Error())
	}

	newUser := model.User{}
	newUser.Email = data.Email
	newUser.Username = data.Username
	newUser.Password = string(hashedPassword)
	newUser.Role = "ADMIN"

	return s.repo.CreateAdmin(newUser)
}

func (s *userService) UserChangePassword(userId1, userId2 int, role string, user *model.UserChangePassword) error {
	fmt.Println(userId1, userId2, role)
	if role == "USER" {

		if userId1 != userId2 {
			return notFound("Permission denied")
		}

		exist, err := s.repo.GetUserByID(userId1)
		if err != nil {
			return internalServerError(err.Error())
		}

		if exist == nil {
			return notFound("User not found")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.OldPassword)); err != nil {
			return notFound("Password not match")
		}

	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError(err.Error())
	}

	user.Password = string(password)

	if err := s.repo.ChangePassword(userId1, user.Password); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *userService) AdminChangePassword(id int, user *model.AdminChangePassword) error {

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return internalServerError(err.Error())
	}

	user.Password = string(password)

	if err := s.repo.ChangePassword(id, user.Password); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}

func (s *userService) DeleteUser(id int) error {

	if err := s.repo.DeleteUser(id); err != nil {
		return internalServerError(err.Error())
	}

	return nil
}
