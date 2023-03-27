package frontend

import (
	"cyber-api/handler"
	"cyber-api/model"
	"cyber-api/service"
	"errors"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type deviceController struct {
	deviceService service.DeviceService
}

func newDeviceController(
	deviceService service.DeviceService,
) deviceController {
	return deviceController{deviceService}
}

func DeviceController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewDeviceRepository(db)
	websiteRepo := repository.NewWebsiteRepository(db)
	service := service.NewDeviceService(repo, websiteRepo)
	handler := newDeviceController(service)

	r = r.Group("/devices")
	r.POST("", handler.createDevice)

}

// @Summary CreateDevice
// @Description create device
// @Tags Devices
// @Accept  json
// @Produce  json
// @Param body body model.DeviceBody true "body"
// @Success 201 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /devices [post]
func (h deviceController) createDevice(c *gin.Context) {

	data := &model.DeviceBody{}

	if err := c.ShouldBindJSON(data); err != nil {
		handler.HandleError(c, err)
		return
	}

	if err := validate(data); err != nil {
		handler.HandleError(c, err)
		return
	}

	token, err := h.deviceService.CreateDevice(data)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithToken{Message: "success", Token: token})
}

func validate(data *model.DeviceBody) error {

	if data.HardwareId == "" {
		return errors.New("HardwareId is required")
	}

	return nil
}
