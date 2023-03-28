package helper

import (
	"cybergame-api/model"
	"os"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(hardwareId string, deviceId int, userId int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hardwareId": hardwareId,
		"deviceId":   deviceId,
		"userId":     userId,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateJWTUser(user *model.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"role":     user.Role,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
