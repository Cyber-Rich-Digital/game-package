package handler

// import (
// 	"cyber-api/handler"
// 	"cyber-api/model"
// 	"cyber-api/repository"
// 	"cyber-api/service"

// 	firebase "firebase.google.com/go"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type messageController struct {
// 	messageService service.MessageService
// }

// func newMessageController(
// 	messageService service.MessageService,
// ) messageController {
// 	return messageController{messageService}
// }

// func MessageController(r *gin.RouterGroup, db *gorm.DB, firebase *firebase.App) {

// 	repo := repository.NewMessageRepository(db)
// 	tagRepo := repository.NewTagRepository(db)
// 	websiteRepo := repository.NewWebsiteRepository(db)
// 	deviceRepo := repository.NewDeviceRepository(db)
// 	notiRepo := repository.NewNotiRepository(db)

// 	service := service.NewMessageService(repo, websiteRepo, tagRepo, deviceRepo, notiRepo, firebase)
// 	handler := newMessageController(service)

// 	r = r.Group("/messages")
// 	r.POST("", handler.pushNoti)

// }

// // @Summary PushNoti
// // @Description push notification
// // @Tags Back Messages
// // @Security BearerAuth
// // @Accept  json
// // @Produce  json
// // @Param body body model.MessageBody true "body"
// // @Success 201 {object} model.Success
// // @Failure 400 {object} handler.ErrorResponse
// // @Router /be/messages [post]
// func (h messageController) pushNoti(c *gin.Context) {

// 	var body model.MessageBody

// 	if err := c.ShouldBindJSON(&body); err != nil {
// 		handler.HandleError(c, err)
// 		return
// 	}

// 	err := h.messageService.CreateMessage(body)
// 	if err != nil {
// 		handler.HandleError(c, err)
// 		return
// 	}

// 	c.JSON(200, model.Success{Message: "success"})
// }
