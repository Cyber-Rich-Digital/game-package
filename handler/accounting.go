package handler

import (
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
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

	r1 := r.Group("/banks")
	r1.GET("/list", handler.getBanks)
	// r.GET("/detail/:id", handler.getBankByCode)
	// r.GET("/code/:id", handler.getBankByCode)

	r2 := r.Group("/bankaccounts")
	r2.GET("/list", handler.getBankAccounts)
	r2.GET("/detail/:id", handler.getBankAccountById)
	r2.POST("", handler.createBankAccount)
	r2.PATCH("/:id", handler.updateBankAccount)
	r2.DELETE("/:id", handler.deleteBankAccount)
}

// @Summary get Bank List
// @Description get all thai Bank List
// @Tags Banks
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.Pagination
// @Router /banks/list [get]
func (h accountingController) getBanks(c *gin.Context) {

	var query model.BankListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBanks(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetBankAccounts
// @Description get BankAccounts
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /bankaccounts/list [get]
func (h accountingController) getBankAccounts(c *gin.Context) {

	var query model.BankAccountListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBankAccounts(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetBankAccount
// @Description get BankAccount by id
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /bankaccounts/detail/{id} [get]
func (h accountingController) getBankAccountById(c *gin.Context) {

	// userId, _ := c.Get("userId")
	// id := int64(userId.(float64))

	var accounting model.BankAccountParam

	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetBankAccountById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateBankAccount
// @Description create new accounting
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param body body model.BankAccountBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /bankaccounts [post]
func (h accountingController) createBankAccount(c *gin.Context) {

	// userId := c.MustGet("userId")
	// toInt := int(userId.(float64))

	var accounting model.BankAccountBody
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateBankAccount(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary UpdateBankAccount
// @Description update accounting
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankAccountBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /bankaccounts/{id} [patch]
func (h accountingController) updateBankAccount(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
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

	if err := h.accountingService.UpdateBankAccount(identifier, body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary DeleteBankAccount
// @Description delete accounting
// @Tags Bank Accounts
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /bankaccounts/{id} [delete]
func (h accountingController) deleteBankAccount(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteBankAccount(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}

	c.JSON(201, model.Success{Message: "Deleted success"})

}
