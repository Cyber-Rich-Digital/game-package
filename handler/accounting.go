package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type accountingController struct {
	accountingService service.AccountingService
}

func newAccountingController(
	accountingService service.AccountingService,
) accountingController {
	return accountingController{accountingService}
}

func AccountingController(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAccountingRepository(db)
	service := service.NewAccountingService(repo)
	handler := newAccountingController(service)

	r = r.Group("/accountings")
	r.GET("/list", middleware.Authorize, handler.getAccountings)
	r.GET("/detail/:id", middleware.Authorize, handler.getAccounting)
	r.POST("", middleware.Authorize, handler.createAccounting)
	r.PATCH("/:id", middleware.Authorize, handler.updateAccounting)
	r.DELETE("/:id", middleware.Authorize, handler.deleteAccounting)
}

// @Summary GetAccountings
// @Description get Accountings
// @Tags Back Accountings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sort query int false "sort"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/accountings/list [get]
func (h accountingController) getAccountings(c *gin.Context) {

	userId := c.MustGet("userId").(float64)
	role := c.MustGet("role").(string)

	var accounting model.BankAccountQuery

	if err := c.ShouldBind(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	accounting.UserId = int(userId)
	accounting.Role = role

	data, err := h.accountingService.GetBankAccounts(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetAccounting
// @Description get Accounting by id
// @Tags Back Accountings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/accountings/detail/{id} [get]
func (h accountingController) getAccounting(c *gin.Context) {

	userId, _ := c.Get("userId")
	id := int(userId.(float64))

	var accounting model.BankAccountParam

	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBankAccountAndTags(accounting, id)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateAccounting
// @Description create new accounting
// @Tags Back Accountings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param body body model.AccountingBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/accountings [post]
func (h accountingController) createAccounting(c *gin.Context) {

	userId := c.MustGet("userId")
	toInt := int(userId.(float64))

	var accounting model.BankAccountBody

	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateBankAccount(accounting, toInt)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary DeleteAccounting
// @Description delete accounting
// @Tags Back Accountings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/accountings/{id} [delete]
func (h accountingController) deleteAccounting(c *gin.Context) {

	id := c.Param("id")

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	toInt, err := strconv.Atoi(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteBankAccount(toInt)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})

}

// @Summary UpdateAccounting
// @Description update accounting
// @Tags Back Accountings
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Param body body model.AccountingBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /be/accountings/{id} [patch]
func (h accountingController) updateAccounting(c *gin.Context) {

	id := c.Param("id")

	if c.MustGet("role").(string) == "USER" {
		HandleError(c, errors.New("Permission denied"))
		return
	}

	// toInt, err := strconv.Atoi(id)
	toPrimaryKey, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	body := model.BankAccountBody{}

	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}

	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.UpdateBankAccount(toPrimaryKey, body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}
