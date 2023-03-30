package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type groupController struct {
	groupService service.GroupService
}

func newGroupController(
	groupService service.GroupService,
) groupController {
	return groupController{groupService}
}

func GroupController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewGroupRepository(db)
	service := service.NewGroupService(repo)
	handler := newGroupController(service)

	r = r.Group("/groups")
	r.POST("/create", handler.create)

}

// @Summary Create Group
// @Description Create Group
// @Tags Groups
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param register body model.CreateGroup true "Create Group"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /groups/create [post]
func (h groupController) create(c *gin.Context) {

	data := &model.CreateGroup{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	err := h.groupService.Create(data)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}
