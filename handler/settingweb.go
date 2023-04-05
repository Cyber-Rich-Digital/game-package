package handler

import (
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type settingwebController struct {
	settingebService service.SettingWebService
}

func newSettingwebController(
	settingebService service.SettingWebService,
) settingwebController {
	return settingwebController{settingebService}
}

func SettingwebController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewSettingWebRepository(db)
	service := service.NewSettingWebService(repo)
	handler := newSettingwebController(service)

	r = r.Group("/settingweb")
	r.POST("/create", handler.create)

}
func (h settingwebController) create(c *gin.Context) {

	data := &model.SettingwebCreateBody{}
	if err := c.ShouldBindJSON(data); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created settingweb success"})
}
