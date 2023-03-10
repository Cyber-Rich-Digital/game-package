package frontend

import (
	"cyber-api/handler"
	"cyber-api/model"
	"cyber-api/service"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
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
// @Tags Front Auth
// @Accept  json
// @Produce  json
// @Param register body model.CreateUser true "Register"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /register [post]
func (h authController) register(c *gin.Context) {

	data := &model.CreateUser{}
	if err := c.ShouldBindJSON(data); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := handler.ValidateField(data); err != nil {
		handler.HandleError(c, err)
		return
	}

	err := h.userService.CreateUser(data)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.Success{Message: "Register success"})
}

// @Summary Login
// @Description Login
// @Tags Front Auth
// @Accept  json
// @Produce  json
// @Param login body model.Login true "Login"
// @Success 201 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /login [post]
func (h authController) login(c *gin.Context) {

	body := &model.Login{}
	if err := c.ShouldBindJSON(body); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := handler.ValidateField(body); err != nil {
		handler.HandleError(c, err)
		return
	}

	token, err := h.userService.Login(body)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithToken{Message: "Login success", Token: token})
}
