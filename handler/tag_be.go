package handler

import (
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/repository"
	"cyber-api/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tagController struct {
	tagService service.TagService
}

func newTagController(
	tagService service.TagService,
) tagController {
	return tagController{tagService}
}

func TagController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewTagRepository(db)
	service := service.NewTagService(repo)
	handler := newTagController(service)

	r = r.Group("/tags")
	r.DELETE("/:id", middleware.Authorize, handler.deleteTag)
}

// @Summary DeleteTag
// @Description delete tag by id
// @Tags Back Tags
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/tags/{id} [delete]
func (h tagController) deleteTag(c *gin.Context) {

	id := c.Param("id")

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.tagService.DeleteTag(toInt); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})
}
