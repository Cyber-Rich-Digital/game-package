package handler

import (
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/service"
	"errors"
	"strconv"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type userController struct {
	userService service.UserService
}

func newUserController(
	userService service.UserService,
) userController {
	return userController{userService}
}

func UserController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewUserRepository(db)
	websiteRepo := repository.NewWebsiteRepository(db)
	service := service.NewUserService(repo, websiteRepo)
	handler := newUserController(service)

	r = r.Group("/users")
	r.GET("/user", middleware.Authorize, handler.getAllUser)
	r.GET("/admin", middleware.Authorize, handler.getAllAdmin)
	r.POST("/create/admin", middleware.Authorize, handler.createAdmin)
	r.PUT("/changepass/user/:id", middleware.Authorize, handler.userChangePassword)
	r.PUT("/changepass/admin/:id", middleware.Authorize, handler.adminChangePassword)
	r.DELETE("/delete/:id", middleware.Authorize, handler.deleteUser)

}

// @Summary Get All User
// @Description Get All User
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sort query int false "sort"
// @Success 200 {object} model.Pagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/user [get]
func (h userController) getAllUser(c *gin.Context) {

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	query := model.UserQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	users, err := h.userService.GetUsers(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, users)
}

// @Summary Get All Admin
// @Description Get All Admin
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sort query int false "sort"
// @Success 200 {object} model.Pagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/admin [get]
func (h userController) getAllAdmin(c *gin.Context) {

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	query := model.UserQuery{}

	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	users, err := h.userService.GetAdmins(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, users)
}

// // @Summary Register
// // @Description Register
// // @Tags Back Auth
// // @Security BearerAuth
// // @Accept  json
// // @Produce  json
// // @Param register body model.CreateUser true "Register"
// // @Success 201 {object} model.Success
// // @Failure 400 {object} handler.ErrorResponse
// // @Router /be/register [post]
// func (h userController) register(c *gin.Context) {

// 	data := &model.CreateUser{}
// 	if err := c.ShouldBindJSON(data); err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	if err := handler.ValidateFieldUser(data); err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	err := h.userService.CreateUser(data)
// 	if err != nil {
// 		HandleError(c, err)
// 		return
// 	}

// 	c.JSON(200, model.Success{Message: "Register success"})
// }

// @Summary Create Admin
// @Description Create Admin
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param register body model.CreateUser true "Register"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/create/admin [post]
func (h userController) createAdmin(c *gin.Context) {

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	data := &model.CreateAdmin{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.userService.CreateAdmin(data, true)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Register success"})
}

// @Summary Change User Password
// @Description Change User Password
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param changePassword body model.UserChangePassword true "Change Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/changepass/user/{id} [put]
func (h userController) userChangePassword(c *gin.Context) {

	userId2 := c.MustGet("userId").(float64)
	role := c.MustGet("role").(string)
	userId1 := c.Param("id")
	toInt, _ := strconv.Atoi(userId1)

	data := &model.UserChangePassword{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.userService.UserChangePassword(toInt, int(userId2), role, data)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Changed password success"})
}

// @Summary Change Admin Password
// @Description Change Admin Password
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Admin ID"
// @Param changePassword body model.AdminChangePassword true "Change Password"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/changepass/admin/{id} [put]
func (h userController) adminChangePassword(c *gin.Context) {

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		HandleError(c, err)
		return
	}

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	data := &model.AdminChangePassword{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.userService.AdminChangePassword(userId, data)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Changed password success"})
}

// @Summary Delete User
// @Description Delete User
// @Tags Back Users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/users/{id} [delete]
func (h userController) deleteUser(c *gin.Context) {

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.userService.DeleteUser(toInt)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})

}
