package frontend

import (
	"cyber-api/handler"
	"cyber-api/middleware"
	"cyber-api/model"
	"cyber-api/service"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userController struct {
	userService service.UserService
}

func newUserController(
	userService service.UserService,
) userController {
	return userController{userService}
}

func UserController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := newUserController(service)

	r = r.Group("/users")
	r.GET("/:id", middleware.Authorize, handler.getUser)

}

// @Summary get user
// @Description User
// @Tags Front User
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Ok 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /users/{id} [get]
func (h userController) getUser(c *gin.Context) {

	var param model.GetParam

	if err := c.ShouldBindUri(&param); err != nil {
		handler.HandleError(c, err)
		return
	}

	result, err := h.userService.GetUserByID(param.Id)
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, model.OKWithResult{Result: result})
}
