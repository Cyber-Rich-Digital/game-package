package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type errorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidateResponse struct {
	Errors []errorMsg `json:"errors"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte":
		return "Should be less than " + fe.Param()
	case "gte":
		return "Should be greater than " + fe.Param()
	case "min":
		return "Should be greater than " + fe.Param()
	case "max":
		return "Should be less than " + fe.Param()
	}
	return "Unknown error"
}

func HandleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case service.ResponseError:
		c.AbortWithStatusJSON(e.Code, ErrorResponse{Message: e.Message})
	case validator.ValidationErrors:
		list := make([]errorMsg, len(e))
		for i, fe := range e {

			if fe.Field() == "Email" {
				errMessage := "Email is invalid"
				list[i] = errorMsg{fe.Field(), errMessage}
				continue
			}

			list[i] = errorMsg{fe.Field(), getErrorMsg(fe)}
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, ValidateResponse{Errors: list})
	case error:
		status := http.StatusBadRequest
		c.AbortWithStatusJSON(status, ErrorResponse{Message: e.Error()})
	}
}

func ValidateField[T model.CreateAdmin | model.LoginAdmin](data T) error {

	if err := validator.New().Struct(data); err != nil {
		checkType := strings.Split(err.(validator.ValidationErrors).Error(), "'")[3]
		if checkType == "Phone" || checkType == "Password" {
			return errors.New("Phone or Password is invalid")
		} else {
			return errors.New("Invalid data")
		}
	}

	return nil
}
