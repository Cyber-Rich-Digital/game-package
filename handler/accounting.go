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

	root := r.Group("/accounting")
	root.GET("/autocreditflags/list", handler.getAutoCreditFlags)
	root.GET("/autowithdrawflags/list", handler.getAutoWithdrawFlags)
	root.GET("/qrwalletstatuses/list", handler.getQrWalletStatuses)
	root.GET("/accountpriorities/list", handler.getAccountPriorities)
	root.GET("/accountstatuses/list", handler.getAccountStatuses)

	bankRoute := root.Group("/banks")
	bankRoute.GET("/list", handler.getBanks)

	accountTypeRoute := root.Group("/accounttypes")
	accountTypeRoute.GET("/list", handler.getAccountTypes)

	accountRoute := root.Group("/bankaccounts")
	accountRoute.GET("/list", handler.getBankAccounts)
	accountRoute.GET("/detail/:id", handler.getBankAccountById)
	accountRoute.POST("", handler.createBankAccount)
	accountRoute.PATCH("/:id", handler.updateBankAccount)
	accountRoute.DELETE("/:id", handler.deleteBankAccount)

	transactionRoute := root.Group("/transactions")
	transactionRoute.GET("/list", handler.getTransactions)
	transactionRoute.GET("/detail/:id", handler.getTransactionById)
	transactionRoute.POST("", handler.createTransaction)
	transactionRoute.DELETE("/:id", handler.deleteTransaction)

	transferRoute := root.Group("/transfers")
	transferRoute.GET("/list", handler.getTransfers)
	transferRoute.GET("/detail/:id", handler.getTransferById)
	transferRoute.POST("", handler.createTransfer)
	transferRoute.POST("/confirm/:id", handler.confirmTransfer)
	transferRoute.DELETE("/:id", handler.deleteTransfer)
}

// @Summary get Bank List
// @Description get all thai Bank List
// @Tags Options
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.Pagination
// @Router /accounting/banks/list [get]
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

// @Summary get Account Type List
// @Description get all Account Type
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/accounttypes/list [get]
func (h accountingController) getAccountTypes(c *gin.Context) {

	var query model.AccountTypeListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetAccountTypes(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary get Auto Credit Flags
// @Description get all Auto Credit Flags
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/autocreditflags/list [get]
func (h accountingController) getAutoCreditFlags(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "manual", Name: "สร้างใบงานและปรับเครดิตเอง"},
		{Key: "auto", Name: "ปรับเครดิตออโต้ (Bot)"},
	}
	c.JSON(200, model.Pagination{List: data, Total: 2})
}

// @Summary get Auto withdraw Flags
// @Description get all Auto withdraw Flags Flags
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/autowithdrawflags/list [get]
func (h accountingController) getAutoWithdrawFlags(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "manual", Name: "สร้างใบงานและปรับเครดิตเอง"},
		{Key: "auto_backoffice", Name: "บัญชีถอนหลัก ปรับเครดิตออโต้ คลิกผ่านระบบหลังบ้าน"},
		{Key: "auto_bot", Name: "บัญชีถอนหลัก ปรับเครดิตออโต้ โอนเงินออโต้ (Bot)"},
	}
	c.JSON(200, model.Pagination{List: data, Total: 2})
}

// @Summary get Qr Wallet Statuses
// @Description get all Qr Wallet Statuses Flags
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/qrwalletstatuses/list [get]
func (h accountingController) getQrWalletStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "use_qr", Name: "เปิด"},
		{Key: "disabled", Name: "ปิด"},
	}
	c.JSON(200, model.Pagination{List: data, Total: 2})
}

// @Summary get Account Statuses
// @Description get all Account Statuses Flags
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/accountstatuses/list [get]
func (h accountingController) getAccountStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "active", Name: "ใช้งาน"},
		{Key: "deactive", Name: "ระงับการใช้งาน"},
	}
	c.JSON(200, model.Pagination{List: data, Total: 2})
}

// @Summary get Account Priorities
// @Description get all Account Priorities Flags
// @Tags Options
// @Accept json
// @Produce json
// @Success 200 {object} model.Pagination
// @Router /accounting/accountpriorities/list [get]
func (h accountingController) getAccountPriorities(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "new", Name: "ระดับ NEW ทั่วไป"},
		{Key: "gold", Name: "ระดับ Gold ฝากมากกว่า 10 ครั้ง"},
		{Key: "platinum", Name: "ระดับ Platinum ฝากมากกว่า 20 ครั้ง"},
		{Key: "vip", Name: "ระดับ VIP ฝากมากกว่า 20 ครั้ง"},
		{Key: "classic", Name: "ระดับ CLASSIC ฝากสะสมมากกว่า 1,000 บาท"},
		{Key: "superior", Name: "ระดับ SUPERIOR ฝากสะสมมากกว่า 10,000 บาท"},
		{Key: "deluxe", Name: "ระดับ DELUXE ฝากสะสมมากกว่า 100,000 บาท"},
		{Key: "wisdom", Name: "ระดับ WISDOM ฝากสะสมมากกว่า 500,000 บาท"},
	}
	c.JSON(200, model.Pagination{List: data, Total: 2})
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
// @Router /accounting/bankaccounts/list [get]
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
// @Router /accounting/bankaccounts/detail/{id} [get]
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
// @Router /accounting/bankaccounts [post]
func (h accountingController) createBankAccount(c *gin.Context) {

	// bankId := c.MustGet("bankId")
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
// @Router /accounting/bankaccounts/{id} [patch]
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
// @Router /accounting/bankaccounts/{id} [delete]
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

// @Summary GetTransactions
// @Description get Transactions
// @Tags Bank Account Transactions
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/list [get]
func (h accountingController) getTransactions(c *gin.Context) {

	var query model.BankAccountTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetTransaction
// @Description get Transaction by id
// @Tags Bank Account Transactions
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/detail/{id} [get]
func (h accountingController) getTransactionById(c *gin.Context) {

	// userId, _ := c.Get("userId")
	// id := int64(userId.(float64))

	var accounting model.BankAccountTransactionParam

	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransactionById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateTransaction
// @Description create new accounting
// @Tags Bank Account Transactions
// @Accept json
// @Produce json
// @Param body body model.BankAccountTransactionBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions [post]
func (h accountingController) createTransaction(c *gin.Context) {

	// bankId := c.MustGet("bankId")
	// toInt := int(userId.(float64))

	var accounting model.BankAccountTransactionBody
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateTransaction(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary DeleteTransaction
// @Description delete accounting
// @Tags Bank Account Transactions
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transactions/{id} [delete]
func (h accountingController) deleteTransaction(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteTransaction(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary GetTransfers
// @Description get Transfers
// @Tags Bank Account Transfers
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/list [get]
func (h accountingController) getTransfers(c *gin.Context) {

	var query model.BankAccountTransferListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransfers(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.Pagination{List: data.List, Total: data.Total})
}

// @Summary GetTransfer
// @Description get Transfer by id
// @Tags Bank Account Transfers
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithToken
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/detail/{id} [get]
func (h accountingController) getTransferById(c *gin.Context) {

	// userId, _ := c.Get("userId")
	// id := int64(userId.(float64))

	var accounting model.BankAccountTransferParam

	if err := c.ShouldBindUri(&accounting); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.accountingService.GetTransferById(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateTransfer
// @Description create new Transfer
// @Tags Bank Account Transfers
// @Accept json
// @Produce json
// @Param body body model.BankAccountTransferBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers [post]
func (h accountingController) createTransfer(c *gin.Context) {

	// bankId := c.MustGet("bankId")
	// toInt := int(userId.(float64))

	var accounting model.BankAccountTransferBody
	if err := c.ShouldBindJSON(&accounting); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(accounting); err != nil {
		HandleError(c, err)
		return
	}

	err := h.accountingService.CreateTransfer(accounting)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary ConfirmTransfer
// @Description update Transfer
// @Tags Bank Account Transfers
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankAccountTransferConfirmBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/confirm/{id} [post]
func (h accountingController) confirmTransfer(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	body := model.BankAccountTransferConfirmBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(body); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.accountingService.ConfirmTransfer(identifier, body); err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(201, model.Success{Message: "Updated success"})
}

// @Summary DeleteTransfer
// @Description delete Transfer
// @Tags Bank Account Transfers
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /accounting/transfers/{id} [delete]
func (h accountingController) deleteTransfer(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.accountingService.DeleteTransfer(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}
