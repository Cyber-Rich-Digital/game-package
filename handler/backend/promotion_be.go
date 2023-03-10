package backend

import (
	"cyber-api/handler"
	"cyber-api/service"

	"cyber-api/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type promotionController struct {
	promotionService service.PromotionService
}

func newPromotionController(
	promotionService service.PromotionService,
) promotionController {
	return promotionController{promotionService}
}

func PromotionController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewPromotionRepository(db)
	service := service.NewPromotionService(repo)
	handler := newPromotionController(service)

	r = r.Group("/promotions")

	r.GET("/", handler.getAll)

}

// @Summary Get all promotions
// @Description Get all promotions
// @Tags Back Promotion
// @Accept  json
// @Produce  json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Success 200 {object} model.Pagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /api/be/promotions [get]
func (h promotionController) getAll(c *gin.Context) {

	result, err := h.promotionService.GetPromotions()
	if err != nil {
		handler.HandleError(c, err)
		return
	}

	c.JSON(200, result)
}
