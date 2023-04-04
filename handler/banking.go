package handler

import (
	"cybergame-api/middleware"
	"cybergame-api/model"
	"cybergame-api/repository"
	"cybergame-api/service"
	"strconv"

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
	transactionRoute.DELETE("/:id", middleware.Authorize, handler.deleteBankTransaction)

}

// @Summary get Transaction Type List
// @Description ดึงข้อมูลตัวเลือก ประเภทรายการ
// @Tags Banking - Options
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} model.SuccessWithPagination
// @Router /banking/transactiontypes/list [get]
func (h bankingController) getTransactionTypes(c *gin.Context) {

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

	c.JSON(200, model.SuccessWithPagination{List: data.List, Total: data.Total})
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
// @Description สร้างข้อมูลการฝากถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.BankTransactionCreateBody true "body"
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

// @Summary DeleteTransaction
// @Description ลบข้อมูลการฝากถอน
// @Tags Banking - Bank Transaction
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/transactions/{id} [delete]
func (h bankingController) deleteBankTransaction(c *gin.Context) {

	id := c.Param("id")
	identifier, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.bankingService.DeleteBankTransaction(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}
