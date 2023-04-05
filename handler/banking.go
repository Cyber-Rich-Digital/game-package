package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type bankingController struct {
	bankingService    service.BankingService
	accountingService service.AccountingService
}

func newBankingController(
	bankingService service.BankingService,
	accountingService service.AccountingService,
) bankingController {
	return bankingController{bankingService, accountingService}
}

func BankingController(r *gin.RouterGroup, db *gorm.DB) {

	repoBanking := repository.NewBankingRepository(db)
	repoAccounting := repository.NewAccountingRepository(db)
	service1 := service.NewBankingService(repoBanking, repoAccounting)
	service2 := service.NewAccountingService(repoAccounting)
	handler := newBankingController(service1, service2)

	root := r.Group("/banking")
	root.GET("/transactiontypes/list", middleware.Authorize, handler.getTransactionTypes)

	statementRoute := root.Group("/statements")
	statementRoute.GET("/list", middleware.Authorize, handler.getBankStatements)
	statementRoute.GET("/detail/:id", middleware.Authorize, handler.getBankStatementById)
	statementRoute.POST("", middleware.Authorize, handler.createBankStatement)
	statementRoute.DELETE("/:id", middleware.Authorize, handler.deleteBankStatement)

	transactionRoute := root.Group("/transactions")
	transactionRoute.GET("/list", middleware.Authorize, handler.getBankTransactions)
	transactionRoute.GET("/detail/:id", middleware.Authorize, handler.getBankTransactionById)
	transactionRoute.POST("", middleware.Authorize, handler.createBankTransaction)
	transactionRoute.POST("/bonus", middleware.Authorize, handler.createBonusTransaction)

	transactionRoute.GET("/pendingdepositlist", middleware.Authorize, handler.getPendingDepositTransactions)
	transactionRoute.GET("/pendingwithdrawlist", middleware.Authorize, handler.getPendingWithdrawTransactions)
	transactionRoute.GET("/finishedlist", middleware.Authorize, handler.getFinishedTransactions)
	transactionRoute.POST("/remove/:id", middleware.Authorize, handler.removeFinishedTransaction)
}

// @Summary get Transaction Type List
// @Description ดึงข้อมูลตัวเลือก ประเภทการทำรายการ
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/transactiontypes/list [get]
func (h bankingController) getTransactionTypes(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "deposit", Name: "ฝาก"},
		{Key: "withdraw", Name: "ถอน"},
		{Key: "getcreditback", Name: "ดึงเครดิตกลับ"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary GetStatementList
// @Description ดึงข้อมูลลิสการโอนเงิน ใช้แสดงในหน้า จัดการธนาคาร - ธุรกรรม
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankStatementListRequest true "BankStatementListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/list [get]
func (h bankingController) getBankStatements(c *gin.Context) {

	var query model.BankStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankStatements(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetStatementByID
// @Description ดึงข้อมูลการโอนด้วย id *ยังไม่ได้ใช้งาน*
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/detail/{id} [get]
func (h bankingController) getBankStatementById(c *gin.Context) {

	var req model.BankStatementGetRequest

	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankStatementById(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateStatement
// @Description สร้างข้อมูลการโอน
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankStatementCreateBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements [post]
func (h bankingController) createBankStatement(c *gin.Context) {

	var banking model.BankStatementCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.CreateBankStatement(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary DeleteStatement
// @Description ลบข้อมูลการโอนด้วย id ใช้ในหน้า จัดการธนาคาร - ธุรกรรม ส่งรหัสผ่านมาเพื่อยืนยันด้วย
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/{id} [delete]
func (h bankingController) deleteBankStatement(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.bankingService.DeleteBankStatement(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}

// @Summary GetTransactionList
// @Description ดึงข้อมูลลิสการฝากถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.BankTransactionListRequest true "BankTransactionListRequest"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/list [get]
func (h bankingController) getBankTransactions(c *gin.Context) {

	var query model.BankTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetTransactionByID
// @Description ดึงข้อมูลการฝากถอน ด้วย id
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/detail/{id} [get]
func (h bankingController) getBankTransactionById(c *gin.Context) {

	var req model.BankTransactionGetRequest

	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetBankTransactionById(req)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary CreateTransaction
// @Description สร้างข้อมูล บันทึกรายการฝาก-ถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankTransactionCreateBody true "*บังคับกรอก memberCode และ creditAmount และ transferType"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions [post]
func (h bankingController) createBankTransaction(c *gin.Context) {

	var banking model.BankTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.CreateBankTransaction(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary CreateTransaction
// @Description สร้างข้อมูล บันทึกแจกโบนัส
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BonusTransactionCreateBody true "body description"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/bonus [post]
func (h bankingController) createBonusTransaction(c *gin.Context) {

	var banking model.BonusTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.CreateBonusTransaction(banking); err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(201, model.Success{Message: "Created success"})
}

// @Summary GetPendingDepositTransactions
// @Description ดึงข้อมูลลิสการฝาก ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.PendingDepositTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/pendingdepositlist [get]
func (h bankingController) getPendingDepositTransactions(c *gin.Context) {

	var query model.PendingDepositTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetPendingDepositTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetPendingWithdrawTransactions
// @Description ดึงข้อมูลลิสการถอน ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.PendingWithdrawTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/pendingwithdrawlist [get]
func (h bankingController) getPendingWithdrawTransactions(c *gin.Context) {

	var query model.PendingWithdrawTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetPendingWithdrawTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetFinishedTransactions
// @Description ดึงข้อมูลลิสการถอน ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.FinishedTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/finishedlist [get]
func (h bankingController) getFinishedTransactions(c *gin.Context) {

	var query model.FinishedTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetFinishedTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary RemoveFinishedTransaction
// @Description ลบข้อมูลการฝากถอนเสร็จสิ้น
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/remove/{id} [post]
func (h bankingController) removeFinishedTransaction(c *gin.Context) {

	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}
	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	var data model.BankTransactionRemoveBody
	data.Status = "removed"
	data.RemovedAt = time.Now()
	data.RemovedByUserId = *adminId
	data.RemovedByUsername = *username

	delErr := h.bankingService.RemoveFinishedTransaction(identifier, data)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}
