package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type authError struct {
	Message string `json:"message" example:"error" `
}

func Authorize(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if len(strings.Split(token, " ")) != 2 {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	token = strings.Split(token, " ")[1]

	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if !claims.Valid {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("id", claims.Claims.(jwt.MapClaims)["id"])

	c.Next()
}

func AdminAuthorize(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if len(strings.Split(token, " ")) != 2 {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	token = strings.Split(token, " ")[1]

	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if !claims.Valid {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("id", claims.Claims.(jwt.MapClaims)["id"])
	c.Set("role", claims.Claims.(jwt.MapClaims)["role"])

	c.Next()
}
