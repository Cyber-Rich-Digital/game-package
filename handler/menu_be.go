package handler

import (
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/service"

	"github.com/gin-gonic/gin"
)

type menuController struct {
	menuService service.MenuService
}

func newMenuController(
	menuService service.MenuService,
) menuController {
	return menuController{menuService}
}

func MenuController(r *gin.RouterGroup) {

	service := service.NewMenuService()
	handler := newMenuController(service)

	r = r.Group("/menu")
	r.GET("", middleware.Authorize, handler.getMenu)
}

// @Summary GetMenu
// @Description get menu by role
// @Tags Back Menu
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ResponseAsList
// @Router /be/menu [get]
func (h menuController) getMenu(c *gin.Context) {

	role := c.MustGet("role").(string)
	result := h.menuService.GetMenu(role)

	c.JSON(200, model.ResponseAsList{List: result})
}
