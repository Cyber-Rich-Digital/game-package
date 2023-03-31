package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"fmt"
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
	r.PUT("/group/:id", middleware.Authorize, handler.updateGroup)
	r.DELETE("/group/:id", middleware.Authorize, handler.deleteGroup)
	r.DELETE("/permission/:id", middleware.Authorize, handler.DeletePermission)

}

// @Summary Group List
// @Description Group List
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group [get]
func (h adminController) groupList(c *gin.Context) {

	query := model.AdminGroupQuery{}
	if err := c.ShouldBindQuery(&query); err != nil {
		HandleError(c, err)
		return
	}

	fmt.Println(query.Page, query.Limit)

	data, err := h.adminService.GetGroupList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, data)
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

// @Summary Update Group
// @Description Update Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Param register body model.AdminUpdateGroup true "Update Group"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [put]
func (h adminController) updateGroup(c *gin.Context) {

	data := &model.AdminUpdateGroup{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	data.GroupId = int64(toInt)

	err = h.adminService.UpdateGroup(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary Delete Group
// @Description Delete Group
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Group ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/group/{id} [delete]
func (h adminController) deleteGroup(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	var param model.DeleteGroup
	param.Id = int64(toInt)

	if err := h.adminService.DeleteGroup(param.Id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary Delete Permission
// @Description Delete Permission
// @Tags Admins
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Permission ID"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /admins/permission/{id} [delete]
func (h adminController) DeletePermission(c *gin.Context) {

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	var param model.DeletePermission
	param.Id = int64(toInt)

	if err := h.adminService.DeletePermission(param.Id); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}
