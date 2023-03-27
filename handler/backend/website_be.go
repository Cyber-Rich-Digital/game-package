package backend

import (
	"cyber-api/handler"
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/repository"
	"cyber-api/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	r.GET("/list", middleware.Authorize, handler.getWebsites)
	r.GET("/detail/:id", middleware.Authorize, handler.getWebsite)
	r.POST("", middleware.Authorize, handler.createWebsite)
	r.PATCH("/:id", middleware.Authorize, handler.updateWebsite)
	r.DELETE("/:id", middleware.Authorize, handler.deleteWebsite)
}

// @Summary GetWebsites
// @Description get Websites
// @Tags Back Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sort query int false "sort"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/websites/list [get]
func (h websiteController) getWebsites(c *gin.Context) {

	userId := c.MustGet("userId").(float64)
	role := c.MustGet("role").(string)

	var website model.WebsiteQuery

	if err := c.ShouldBind(&website); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := validator.New().Struct(website); err != nil {
		handler.HandleError(c, err)
		return
	}

	website.UserId = int(userId)
	website.Role = role

	data, err := h.websiteService.GetWebsites(website)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetWebsite
// @Description get Website by id
// @Tags Back Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/websites/detail/{id} [get]
func (h websiteController) getWebsite(c *gin.Context) {

	userId, _ := c.Get("userId")
	id := int(userId.(float64))

	var website model.WebsiteParam

	if err := c.ShouldBindUri(&website); err != nil {
		handler.HandleError(c, err)
		return
	}

	data, err := h.websiteService.GetWebsiteAndTags(website, id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateWebsite
// @Description create new website
// @Tags Back Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param body body model.WebsiteBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/websites [post]
func (h websiteController) createWebsite(c *gin.Context) {

	userId := c.MustGet("userId")
	toInt := int(userId.(float64))

	var website model.WebsiteBody

	if err := c.ShouldBindJSON(&website); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := validator.New().Struct(website); err != nil {
		handler.HandleError(c, err)
		return
	}

	err := h.websiteService.CreateWebsite(website, toInt)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary DeleteWebsite
// @Description delete website
// @Tags Back Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/websites/{id} [delete]
func (h websiteController) deleteWebsite(c *gin.Context) {

	id := c.Param("id")

	if c.MustGet("role").(string) == "USER" {
		handler.HandleError(c, errors.New("Permission denied"))
		return
	}

	toInt, err := strconv.Atoi(id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	delErr := h.websiteService.DeleteWebsite(toInt)
	if delErr != nil {
		handler.HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})

}

// @Summary UpdateWebsite
// @Description update website
// @Tags Back Websites
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Param body body model.WebsiteBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/websites/{id} [patch]
func (h websiteController) updateWebsite(c *gin.Context) {

	id := c.Param("id")

	if c.MustGet("role").(string) == "USER" {
		handler.HandleError(c, errors.New("Permission denied"))
		return
	}

	toInt, err := strconv.Atoi(id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	body := model.WebsiteBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := h.websiteService.UpdateWebsite(toInt, body); err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}
