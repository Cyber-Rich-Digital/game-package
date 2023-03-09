package service

import (
	"cyber-api/model"
	"cyber-api/repository"
	"errors"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var LoginFailed = errors.New("Email Or Password is incorrect")

type UserService interface {
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
		return nil, err
	}

	if exist.ID == 0 {
		return nil, LoginFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(user.Password)); err != nil {
		return nil, LoginFailed
	}

	token, err := createJWT(&exist)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (s *userService) CreateUser(user *model.CreateUser) error {

	exist, err := s.repo.CheckUserByEmail(user.Email)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("Email already exist")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return s.repo.CreateUser(user)
}

func createJWT(user *model.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
