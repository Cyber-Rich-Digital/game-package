package helper

import (
	"cybergame-api/model"
	"os"

	"github.com/golang-jwt/jwt"
)

func CreateJWT(data model.Admin) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"adminId":  data.Id,
		"role":     data.Role,
		"phone":    data.Phone,
		"username": data.Username,
		"email":    data.Email,
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
