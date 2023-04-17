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
	root.GET("/transactionstatuses/list", middleware.Authorize, handler.getTransactionStatuses)
	root.GET("/statementtypes/list", middleware.Authorize, handler.getStatementTypes)

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
	transactionRoute.POST("/cancel/:id", middleware.Authorize, handler.cancelPendingTransaction)
	transactionRoute.POST("/confirmdeposit/:id", middleware.Authorize, handler.confirmDepositTransaction)
	transactionRoute.POST("/confirmwithdraw/:id", middleware.Authorize, handler.confirmWithdrawTransaction)
	transactionRoute.GET("/finishedlist", middleware.Authorize, handler.getFinishedTransactions)
	transactionRoute.POST("/remove/:id", middleware.Authorize, handler.removeFinishedTransaction)
	transactionRoute.GET("/removedlist", middleware.Authorize, handler.getRemovedTransactions)

	memberRoute := root.Group("/member")
	memberRoute.GET("/transactionsummary", middleware.Authorize, handler.getMemberTransactionSummary)
	memberRoute.GET("/transactions", middleware.Authorize, handler.getMemberTransactions)

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

// @Summary get Transaction Status List
// @Description ดึงข้อมูลตัวเลือก สถานะรายการฝากถอน
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/transactionstatuses/list [get]
func (h bankingController) getTransactionStatuses(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "pending", Name: "รอดำเนินการ"},
		{Key: "canceled", Name: "ยกเลิกแล้ว"},
		{Key: "finished", Name: "อนุมัติแล้ว"},
	}
	c.JSON(200, model.SuccessWithPagination{List: data, Total: 2})
}

// @Summary get Statement type List
// @Description ดึงข้อมูลตัวเลือก ประเภทรายการเดินบัญชี
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/statementtypes/list [get]
func (h bankingController) getStatementTypes(c *gin.Context) {
	var data = []model.SimpleOption{
		{Key: "transfer_in", Name: "โอนเข้า"},
		{Key: "transfer_out", Name: "โอนออก"},
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
// @Success 200 {object} model.SuccessWithPagination
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

	actionErr := h.bankingService.DeleteBankStatement(identifier)
	if actionErr != nil {
		HandleError(c, actionErr)
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
// @Success 200 {object} model.SuccessWithPagination
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

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
	if err != nil {
		HandleError(c, err)
		return
	}

	var banking model.BankTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}
	banking.CreatedByUserId = *adminId
	banking.CreatedByUsername = *username

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

	var banking model.BonusTransactionCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}
	banking.CreatedByUserId = *adminId
	banking.CreatedByUsername = *username

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
// @Success 200 {object} model.SuccessWithPagination
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
// @Success 200 {object} model.SuccessWithPagination
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

// @Summary CancelPendingTransaction
// @Description ยกเลิก ไม่ยืนยัน ข้อมูลการฝากถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankTransactionCancelBody true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/cancel/{id} [post]
func (h bankingController) cancelPendingTransaction(c *gin.Context) {

	adminId, err := h.accountingService.CheckCurrentAdminId(c.MustGet("adminId"))
	if err != nil {
		HandleError(c, err)
		return
	}
	username, err := h.accountingService.CheckCurrentUsername(c.MustGet("username"))
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

	var data model.BankTransactionCancelBody
	if err := c.ShouldBind(&data); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}

	data.Status = "canceled"
	// data.CancelRemark = data.CancelRemark
	data.CanceledAt = time.Now()
	data.CanceledByUserId = *adminId
	data.CanceledByUsername = *username

	actionErr := h.bankingService.CancelPendingTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Cancel success"})
}

// @Summary ConfirmDepositTransaction
// @Description ยืนยัน ข้อมูลการฝาก
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmDepositRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmdeposit/{id} [post]
func (h bankingController) confirmDepositTransaction(c *gin.Context) {

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

	var data model.BankConfirmDepositRequest
	if err := c.ShouldBind(&data); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}
	data.ConfirmedAt = time.Now()
	data.ConfirmedByUserId = *adminId
	data.ConfirmedByUsername = *username

	actionErr := h.bankingService.ConfirmDepositTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary ConfirmWithdrawTransaction
// @Description ยืนยัน ข้อมูลการถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param body body model.BankConfirmWithdrawRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/confirmwithdraw/{id} [post]
func (h bankingController) confirmWithdrawTransaction(c *gin.Context) {

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

	var data model.BankConfirmWithdrawRequest
	if err := c.ShouldBind(&data); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(data); err != nil {
		HandleError(c, err)
		return
	}
	data.ConfirmedAt = time.Now()
	data.ConfirmedByUserId = *adminId
	data.ConfirmedByUsername = *username

	actionErr := h.bankingService.ConfirmWithdrawTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Confirm success"})
}

// @Summary GetFinishedTransactions
// @Description ดึงข้อมูลลิสการถอน ที่รออนุมัติ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.FinishedTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
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
// @Description ลบข้อมูลการฝากถอน ที่เสร็จสิ้นไปแล้ว
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

	actionErr := h.bankingService.RemoveFinishedTransaction(identifier, data)
	if actionErr != nil {
		HandleError(c, actionErr)
		return
	}
	c.JSON(201, model.Success{Message: "Remove success"})
}

// @Summary GetRemovedTransactions
// @Description ดึงข้อมูลลิสการฝากถอนที่ถูกลบ
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.RemovedTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/removedlist [get]
func (h bankingController) getRemovedTransactions(c *gin.Context) {

	var query model.RemovedTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetRemovedTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}

// @Summary GetMemberTransactionSummary
// @Description ดึงข้อมูลสรุปการฝากถอนของสมาชิก
// @Tags Banking - Member Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/transactionsummary [get]
func (h bankingController) getMemberTransactionSummary(c *gin.Context) {

	var query model.MemberTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberTransactionSummary(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithData{Message: "success", Data: data})
}

// @Summary GetMemberTransactions
// @Description ดึงข้อมูลลิสการฝากถอนของสมาชิก
// @Tags Banking - Member Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param _ query model.MemberTransactionListRequest true "query"
// @Success 200 {object} model.SuccessWithPagination
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/member/transactions [get]
func (h bankingController) getMemberTransactions(c *gin.Context) {

	var query model.MemberTransactionListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetMemberTransactions(query)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
}
