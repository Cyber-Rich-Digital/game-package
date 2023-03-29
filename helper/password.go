package helper

import (
	"os"

	"golang.org/x/crypto/bcrypt"
)

func GenPassword(password string) (string, error) {

	println(password + os.Getenv("PASSWORD_SECRET"))

	hash, err := bcrypt.GenerateFromPassword([]byte(password+os.Getenv("PASSWORD_SECRET")), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(password string, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+os.Getenv("PASSWORD_SECRET")))
	if err != nil {
		return err
	}

	return nil
}
