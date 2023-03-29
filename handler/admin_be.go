package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"errors"

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
	service := service.NewAdminService(repo)
	handler := newUserController(service)

	r = r.Group("/admins")
	r.POST("/create", middleware.Authorize, handler.create)

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

	err := h.adminService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Register success"})
}
