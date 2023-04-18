package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type linenotifyController struct {
	linenotifyService service.LineNotifyService
}

func newLineNotifyController(
	linenotifyService service.LineNotifyService,
) linenotifyController {
	return linenotifyController{linenotifyService}
}

// @Summary CreateLineNotify
// @Description ตั้งค่าแจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.LinenotifyCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/create [post]
func LineNotifyController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewLineNotifyRepository(db)
	service := service.NewLineNotifyService(repo)
	handler := newLineNotifyController(service)

	linenotifRoute := r.Group("/linenotify")
	linenotifRoute.POST("/create", middleware.Authorize, handler.createLineNotify)
	linenotifRoute.GET("/detail/:id", middleware.Authorize, handler.getLineNotifyById)

}
func (h linenotifyController) createLineNotify(c *gin.Context) {

	var line model.LinenotifyCreateBody
	if err := c.ShouldBindJSON(&line); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(line); err != nil {
		HandleError(c, err)
		return
	}

	errline := h.linenotifyService.CreateLineNotify(line)
	if errline != nil {
		HandleError(c, errline)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})

}

// @Summary GetLineNotifyById
// @Description ดึงข้อมูลการcแจ้งเตือนไลน์ ด้วย id
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/detail/{id} [get]
func (h linenotifyController) getLineNotifyById(c *gin.Context) {

	var line model.LinenotifyParam

	if err := c.ShouldBindUri(&line); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.linenotifyService.GetLineNotifyById(line)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}
