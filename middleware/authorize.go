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

	if claims.Claims.(jwt.MapClaims)["deviceId"] != nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	if claims.Claims.(jwt.MapClaims)["id"] == nil && claims.Claims.(jwt.MapClaims)["role"] == nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("userId", claims.Claims.(jwt.MapClaims)["id"])
	c.Set("role", claims.Claims.(jwt.MapClaims)["role"])

	c.Next()
}

func AppAuthorize(c *gin.Context) {

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

	if claims.Claims.(jwt.MapClaims)["id"] == nil &&
		claims.Claims.(jwt.MapClaims)["deviceId"] == nil &&
		claims.Claims.(jwt.MapClaims)["hardwareId"] == nil {
		c.AbortWithStatusJSON(401, authError{
			Message: "Unauthorized",
		})
		return
	}

	c.Set("userId", claims.Claims.(jwt.MapClaims)["userId"])
	c.Set("deviceId", claims.Claims.(jwt.MapClaims)["deviceId"])
	c.Set("hardwareId", claims.Claims.(jwt.MapClaims)["hardwareId"])

	c.Next()
}
