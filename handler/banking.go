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
	// root.GET("/autocreditflags/list", middleware.Authorize, handler.getAutoCreditFlags)

	statementRoute := root.Group("/statements")
	statementRoute.GET("/list", middleware.Authorize, handler.getStatements)
	statementRoute.GET("/detail/:id", middleware.Authorize, handler.getStatementById)
	statementRoute.POST("", middleware.Authorize, handler.createStatement)
	statementRoute.DELETE("/:id", middleware.Authorize, handler.deleteStatement)
}

// @Summary GetStatementList
// @Description ดึงข้อมูลลิสการโอนเงิน ใช้แสดงในหน้า จัดการธนาคาร - ธุรกรรม
// @Tags Banking - Bank Account Statements
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param accountId query string false "accountId"
// @Param amount query string false "amount"
// @Param fromTransferDate query string false "fromTransferDate"
// @Param toTransferDate query string false "toTransferDate"
// @Param status query string false "status"
// @Param search query string false "search"
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param sortCol query string false "sortCol"
// @Param sortAsc query string false "sortAsc"
// @Success 200 {object} model.SuccessWithData
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/list [get]
func (h bankingController) getStatements(c *gin.Context) {

	var query model.BankStatementListRequest
	if err := c.ShouldBind(&query); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(query); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetStatements(query)
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
func (h bankingController) getStatementById(c *gin.Context) {

	var req model.BankStatementGetRequest

	if err := c.ShouldBindUri(&req); err != nil {
		HandleError(c, err)
		return
	}

	data, err := h.bankingService.GetStatementById(req)
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
func (h bankingController) createStatement(c *gin.Context) {

	var banking model.BankStatementCreateBody
	if err := c.ShouldBindJSON(&banking); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(banking); err != nil {
		HandleError(c, err)
		return
	}

	if err := h.bankingService.CreateStatement(banking); err != nil {
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
// @Param body body model.ConfirmRequest true "body"
// @Success 201 {object} model.Success
// @Failure 400 {object} handler.ErrorResponse
// @Router /banking/statements/{id} [delete]
func (h bankingController) deleteStatement(c *gin.Context) {

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

	var confirmation model.ConfirmRequest
	confirmation.UserId = *adminId
	if err := c.ShouldBindJSON(&confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if err := validator.New().Struct(confirmation); err != nil {
		HandleError(c, err)
		return
	}
	if _, err := h.accountingService.CheckConfirmationPassword(confirmation); err != nil {
		HandleError(c, err)
		return
	}

	delErr := h.bankingService.DeleteStatement(identifier)
	if delErr != nil {
		HandleError(c, delErr)
		return
	}
	c.JSON(201, model.Success{Message: "Deleted success"})
}
