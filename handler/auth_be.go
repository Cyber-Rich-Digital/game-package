package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

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
	websiteRepo := repository.NewWebsiteRepository(db)
	service := service.NewUserService(repo, websiteRepo)
	handler := newAuthController(service)

	r.POST("/login", handler.login)
	r.POST("/register", handler.register)
	r.POST("/register/admin", handler.registerAdmin)

}

// @Summary Login
// @Description Login
// @Tags Back Auth
// @Accept  json
// @Produce  json
// @Param login body model.Login true "Login"
// @Success 201 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/login [post]
func (h authController) login(c *gin.Context) {

	body := &model.Login{}
	if err := c.ShouldBindJSON(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := ValidateField(body); err != nil {
		HandleError(c, err)
		return
	}

	token, err := h.userService.Login(body)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.SuccessWithToken{Message: "Login success", Token: token})
}

// @Summary Register
// @Description Register
// @Tags Back Auth
// @Accept  json
// @Produce  json
// @Param register body model.CreateUser true "Register"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/register [post]
func (h authController) register(c *gin.Context) {

	data := &model.CreateUser{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.userService.CreateUser(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Registered success"})
}

// @Summary Register Admin
// @Description Register Admin
// @Tags Back Auth
// @Accept  json
// @Produce  json
// @Param register body model.CreateUser true "Register"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/register/admin [post]
func (h authController) registerAdmin(c *gin.Context) {

	data := &model.CreateAdmin{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := ValidateFieldAdmin(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.userService.CreateAdmin(data, false)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Registered success"})
}
