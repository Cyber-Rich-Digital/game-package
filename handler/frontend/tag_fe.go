package frontend

import (
	"cyber-api/handler"
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/repository"
	"cyber-api/service"

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

	r.GET("/tags/:website_id", middleware.AppAuthorize, handler.getTagByWebsiteId)

}

// @Summary GetTagByWebsiteId
// @Description get all tags by websiteId
// @Tags Tags
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param website_id path int true "website_id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /tags/{website_id} [get]
func (h tagController) getTagByWebsiteId(c *gin.Context) {

	deviceId := c.MustGet("deviceId").(float64)

	var tag model.TagParam

	if err := c.ShouldBindUri(&tag); err != nil {
		handler.HandleError(c, err)
		return
	}

	tag.DeviceId = deviceId

	data, err := h.tagService.GetTagsByWebsiteId(tag)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}
