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

type websiteController struct {
	websiteService service.WebsiteService
}

func newWebsiteController(
	websiteService service.WebsiteService,
) websiteController {
	return websiteController{websiteService}
}

func WebsiteController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewWebsiteRepository(db)
	service := service.NewWebsiteService(repo)
	handler := newWebsiteController(service)

	r = r.Group("/websites")
	r.GET("/:id", middleware.AppAuthorize, handler.getWebsite)
	r.GET("/totals/:date", middleware.AppAuthorize, handler.getWebsiteTotals)
}

// @Summary GetWebsite
// @Description get Website by id
// @Tags Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /websites/{id} [get]
func (h websiteController) getWebsite(c *gin.Context) {

	userId := c.MustGet("userId").(float64)

	var website model.WebsiteParam

	if err := c.ShouldBindUri(&website); err != nil {
		handler.HandleError(c, err)
		return
	}

	data, err := h.websiteService.GetWebsite(website, int(userId))
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetWebsiteTotals
// @Description get Website totals
// @Tags Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param date path string true "date"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /websites/totals/{date} [get]
func (h websiteController) getWebsiteTotals(c *gin.Context) {

	userId := c.MustGet("userId").(float64)

	var body model.WebsiteDate

	if err := c.ShouldBindUri(&body); err != nil {
		handler.HandleError(c, err)
		return
	}

	body.UserId = int(userId)

	data, err := h.websiteService.GetWebsiteTotals(body)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}
