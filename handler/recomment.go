package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/service"
	"strconv"

	"cybergame-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type recommentController struct {
	recommentService service.RecommentService
}

func newRecommentController(
	recommentService service.RecommentService,
) recommentController {
	return recommentController{recommentService}
}

func RecommentController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewRecommentRepository(db)
	service := service.NewRecommentService(repo)
	handler := newRecommentController(service)

	r = r.Group("/recomments")
	r.GET("/list", middleware.Authorize, handler.getRecommentList)
	r.POST("/create", middleware.Authorize, handler.createRecomment)
	r.PUT("/update/:id", middleware.Authorize, handler.updateRecomment)

}

// @Summary Get Recomment List
// @Description Get Recomment List
// @Tags Recomments
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param _ query model.RecommentQuery true "Query Recomment"
// @Success 200 {object} model.SuccessWithList
// @Failure 400 {object} handler.ErrorResponse
// @Router /recomments/list [get]
func (h recommentController) getRecommentList(c *gin.Context) {

	var query model.RecommentQuery
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}

	list, total, err := h.recommentService.GetRecommentList(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{Message: "Success", List: list, Total: total})
}

// @Summary Create Recomment
// @Description Create Recomment
// @Tags Recomments
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param Body body model.CreateRecomment true "Create Recomment"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /recomments/create [post]
func (h recommentController) createRecomment(c *gin.Context) {

	var body model.CreateRecomment
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.recommentService.CreateRecomment(body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary Update Recomment
// @Description Update Recomment
// @Tags Recomments
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "Recomment ID"
// @Param Body body model.CreateRecomment true "Update Recomment"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /recomments/update/{id} [put]
func (h recommentController) updateRecomment(c *gin.Context) {

	var body model.CreateRecomment
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

	if err := h.recommentService.UpdateRecomment(int64(toInt), body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}
