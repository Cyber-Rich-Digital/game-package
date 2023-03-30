package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"errors"
	"strconv"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type adminController struct {
	adminService service.AdminService
}

func newUserController(
	adminService service.AdminService,
) adminController {
	return adminController{adminService}
}

func AdminController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAdminRepository(db)
	perRepo := repository.NewPermissionRepository(db)
	groupRepo := repository.NewGroupRepository(db)
	service := service.NewAdminService(repo, perRepo, groupRepo)
	handler := newUserController(service)

	r = r.Group("/admins")
	r.GET("/group", middleware.Authorize, handler.groupList)
	r.GET("/group/:id", middleware.Authorize, handler.getGroup)
	r.POST("/create", middleware.Authorize, handler.create)
	r.POST("/creategroup", middleware.Authorize, handler.createGroup)

}

// @Summary Group List
// @Description Group List
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group [get]
func (h adminController) groupList(c *gin.Context) {

	data, err := h.adminService.GetGroupList()
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithList{Message: "Success", List: data})
}

// @Summary Get Group
// @Description Get Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [get]
func (h adminController) getGroup(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		HandleError(c, errors.New("id is required"))
		return
	}

	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.adminService.GetGroup(toInt)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "Success", Data: data})
}

// @Summary Create Admin
// @Description Create Admin
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param register body model.CreateAdmin true "Create Admin"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/create [post]
func (h adminController) create(c *gin.Context) {

	data := &model.CreateAdmin{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.adminService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Register success"})
}

// @Summary Create Group
// @Description Create Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param register body model.AdminCreateGroup true "Create Group"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/creategroup [post]
func (h adminController) createGroup(c *gin.Context) {

	data := &model.AdminCreateGroup{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.adminService.CreateGroup(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}
