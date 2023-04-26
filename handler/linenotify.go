package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"fmt"
	"strconv"

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
	linenotifRoute.PUT("/update/:id", middleware.Authorize, handler.updateLineNotify)

	//GameCyberNoitfy
	linenotifRoute.GET("/typegame/detail/:id", handler.GetLineNotifyGameById)
	linenotifRoute.POST("usergame/create", middleware.Authorize, handler.createLineNotifyGame)
	linenotifRoute.GET("/usergame/detail/:id", middleware.Authorize, handler.GetLineNotifyGameUserById)
	linenotifRoute.DELETE("/usergame/:id", middleware.Authorize, handler.DeleteNotifyGameUser)
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

// @Summary UpdateNotify
// @Description แก้ไข แจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.LinenotifyUpdateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/update/{id} [put]
func (h linenotifyController) updateLineNotify(c *gin.Context) {

	var body model.LinenotifyUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	if err := h.linenotifyService.UpdateLineNotify(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary GetLineNotifyGameById
// @Description ดึงข้อมูลการแจ้งเตือนไลน์ ด้วย id
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/typegame/detail/{id} [get]
func (h linenotifyController) GetLineNotifyGameById(c *gin.Context) {

	//var user model.User
	code := c.Param("code")
	access_token := c.Query("access_token")
	var lineGame model.LinenotifyGameParam

	lineGame.Code = code
	lineGame.Token = access_token

	if err := c.ShouldBindUri(&lineGame); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.linenotifyService.GetLineNotifyGameById(lineGame)
	if err != nil {
		HandleError(c, err)
		return
	}
	typeNotify_Id := fmt.Sprintf("%d", lineGame.Id)

	var bot model.LineNoifyUsergameBody
	bot.UserId = 1
	bot.TypeNotifyId = typeNotify_Id
	bot.Token = access_token

	if err := c.ShouldBindJSON(&bot); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(bot); err != nil {
		HandleError(c, err)
		return
	}

	errline := h.linenotifyService.CreateNotifyGame(bot)
	if errline != nil {
		HandleError(c, errline)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})

}

// @Summary CreateLineNotifyGame
// @Description ตั้งค่าแจ้งเตือนไลน์
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.LineNoifyUsergameBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/usergame/create [post]
func (h linenotifyController) createLineNotifyGame(c *gin.Context) {
	var bot model.LineNoifyUsergameBody
	if err := c.ShouldBindJSON(&bot); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(bot); err != nil {
		HandleError(c, err)
		return
	}

	errline := h.linenotifyService.CreateNotifyGame(bot)
	if errline != nil {
		HandleError(c, errline)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})

}

// @Summary GetLineNotifyGameUserById
// @Description ดึงข้อมูลการแจ้งเตือนไลน์ ด้วย id
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/usergame/detail/{id} [get]
func (h linenotifyController) GetLineNotifyGameUserById(c *gin.Context) {

	var botuser model.LineNotifyUserGameParam

	if err := c.ShouldBindUri(&botuser); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.linenotifyService.GetLineNoifyUserGameById(botuser)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary DeleteLineNoifyUserGame
// @Description ลบข้อมูลการ TOKEN ด้วย id user
// @Tags LineNotify
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.LineNoifyUsergameBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /linenotify/usergame/{id} [delete]
func (h linenotifyController) DeleteNotifyGameUser(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}
	delErr := h.linenotifyService.DeleteLineNoifyUserGame(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}
