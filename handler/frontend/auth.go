package frontend

import (
	"cyber-api/model"
	"cyber-api/service"
	"errors"
	"strings"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type authController struct {
	userService service.UserService
}

func newAuthController(
	userService service.UserService,
) authController {
	return authController{userService}
}

func AuthController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := newAuthController(service)

	r.POST("/register", handler.register)
	r.POST("/login", handler.login)

}

// @Summary Register
// @Description Register
// @Tags Front
// @Accept  json
// @Produce  json
// @Param register body model.CreateUser true "Register"
// @Success 201 {object} model.Success
// @Router /register [post]
func (h authController) register(c *gin.Context) {

	data := &model.CreateUser{}
	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	if err := validateField(data); err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	err := h.userService.CreateUser(data)
	if err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	c.JSON(200, model.Success{Message: "Register success"})
}

// @Summary Login
// @Description Login
// @Tags Front
// @Accept  json
// @Produce  json
// @Param login body model.Login true "Login"
// @Success 201 {object} model.SuccessWithToken
// @Failure 400 {object} model.Error
// @Router /auth/login [post]
func (h authController) login(c *gin.Context) {

	body := &model.Login{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	if err := validateField(body); err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	token, err := h.userService.Login(body)
	if err != nil {
		c.JSON(400, model.Error{Message: err.Error()})
		return
	}

	c.JSON(200, model.SuccessWithToken{Message: "Login success", Token: token})
}

func validateField[T any](data T) error {

	if err := validator.New().Struct(data); err != nil {
		checkType := strings.Split(err.(validator.ValidationErrors).Error(), "'")[3]
		if checkType == "Email" || checkType == "Password" {
			return errors.New("Email or Password is invalid")
		} else {
			return errors.New("Invalid data")
		}
	}

	return nil
}
