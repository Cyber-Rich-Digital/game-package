package model

import (
	"time"

	"gorm.io/gorm"
)

type Member struct {
	Id            int64     `json:"id"`
	MemberCode    string    `json:"memberCode"`
	Username      string    `json:"username"`
	Phone         string    `json:"phone"`
	Firstname     string    `json:"firstname"`
	Lastname      string    `json:"lastname"`
	Fullname      string    `json:"fullname"`
	Credit        float64   `json:"credit"`
	Bankname      string    `json:"bankname"`
	BankAccount   string    `json:"bankAccount"`
	Promotion     string    `json:"promotion"`
	Status        string    `json:"status"`
	Channel       string    `json:"channel"`
	TrueWallet    string    `json:"trueWallet"`
	Note          string    `json:"note"`
	TurnoverLimit int       `json:"turnoverLimit"`
	CreatedAt     time.Time `json:"createdAt"`
}

type MemberListRequest struct {
	Search  string `form:"search" extensions:"x-order:1"`
	Page    int    `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit   int    `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol string `form:"sortCol" extensions:"x-order:9"`
	SortAsc string `form:"sortAsc" extensions:"x-order:10"`
}

type MemberPossibleListRequest struct {
	UnknownStatementId int64   `form:"unknownStatementId" extensions:"x-order:1"`
	UserAccountNumber  *string `form:"userAccountNumber" extensions:"x-order:2"`
	UserBankCode       *string `form:"userBankCode" extensions:"x-order:3"`
	Page               int     `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit              int     `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol            string  `form:"sortCol" extensions:"x-order:9"`
	SortAsc            string  `form:"sortAsc" extensions:"x-order:10"`
}

type BankStatement struct {
	Id                int64          `json:"id" gorm:"primaryKey"`
	AccountId         int64          `json:"accountId"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	Detail            string         `json:"detail"`
	BankId            int64          `json:"bankId"`
	StatementType     string         `json:"statementType"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankCode      string         `json:"fromBankCode"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	FromBankName      string         `json:"fromBankName"`
	FromBankIconUrl   string         `json:"fromBankIconUrl"`
	TransferAt        time.Time      `json:"transferAt"`
	Status            string         `json:"status"`
	CreatedAt         time.Time      `json:"createAt"`
	UpdatedAt         *time.Time     `json:"updateAt"`
	DeletedAt         gorm.DeletedAt `json:"deleteAt"`
}

type BankStatementGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankStatementSummary struct {
	TotalPendingStatementCount int64 `json:"totalPendingStatementCount"`
	TotalPendingDepositCount   int64 `json:"totalPendingDepositCount"`
	TotalPendingWithdrawCount  int64 `json:"totalPendingWithdrawCount"`
}

type BankStatementListRequest struct {
	AccountId        string `form:"accountId" extensions:"x-order:1"`
	StatementType    string `form:"statementType" extensions:"x-order:2"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	SimilarId        int64  `form:"similarId" extensions:"x-order:6"`
	Status           string `form:"status" extensions:"x-order:7"`
	Page             int    `form:"page" extensions:"x-order:7" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:8" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:9"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:10"`
}

type BankStatementCreateBody struct {
	Id                int64     `json:"id"`
	AccountId         int64     `json:"accountId"`
	ExternalId        int64     `json:"externalId"`
	Amount            float32   `json:"amount" sql:"type:decimal(14,2);"`
	Detail            string    `json:"detail"`
	FromBankId        int64     `json:"fromBankId"`
	FromAccountNumber string    `json:"fromAccountNumber"`
	StatementType     string    `json:"statementType"`
	TransferAt        time.Time `json:"transferAt"`
	Status            string    `json:"-"`
}

type BankStatementMatchRequest struct {
	UserId              int64     `json:"userId" validate:"required"`
	ConfirmedAt         time.Time `json:"-"`
	ConfirmedByUserId   int64     `json:"-"`
	ConfirmedByUsername string    `json:"-"`
}

type BankStatementUpdateBody struct {
	Status string `json:"status" validate:"required"`
}

type BankStatementResponse struct {
	Id              int64          `json:"id" gorm:"primaryKey"`
	AccountId       int64          `json:"accountId"`
	AccountName     string         `json:"accountName"`
	AccountNumber   string         `json:"accountNumber"`
	BankName        string         `json:"bankName"`
	Amount          float32        `json:"amount" sql:"type:decimal(14,2);"`
	Detail          string         `json:"detail"`
	FromBankId      int64          `json:"fromBankId"`
	FromBankName    string         `json:"fromBankName"`
	FromBankIconUrl string         `json:"fromBankIconUrl"`
	StatementType   string         `json:"statementType"`
	TransferAt      time.Time      `json:"transferAt"`
	Status          string         `json:"status"`
	CreatedAt       time.Time      `json:"createAt"`
	UpdatedAt       *time.Time     `json:"updateAt"`
	DeletedAt       gorm.DeletedAt `json:"deleteAt"`
}

type BankTransaction struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	MemberCode          string         `json:"memberCode"`
	UserId              int64          `json:"userId"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float32        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float32        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float32        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float32        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float32        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float32        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float32        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          time.Time      `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           time.Time      `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	IsAutoCredit        bool           `json:"isAutoCredit"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}

type BankTransactionGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankTransactionListRequest struct {
	MemberCode       string `form:"memberCode" extensions:"x-order:2"`
	UserId           string `form:"userId" extensions:"x-order:3"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:4"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:5"`
	TransferType     string `form:"transferType" extensions:"x-order:6"`
	TransferStatus   string `form:"transferStatus" extensions:"x-order:7"`
	Search           string `form:"search" extensions:"x-order:8"`
	Page             int    `form:"page" extensions:"x-order:9" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:10" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:11"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:12"`
}

type BankTransactionCreateBody struct {
	Id                int64      `json:"-"`
	MemberCode        string     `json:"memberCode" validate:"required"`
	UserId            int64      `json:"-"`
	TransferType      string     `json:"transferType" validate:"required" example:"deposit"`
	PromotionId       *int64     `json:"promotionId"`
	FromAccountId     *int64     `json:"fromAccountId"`
	FromBankId        *int64     `json:"-"`
	FromAccountName   *string    `json:"-"`
	FromAccountNumber *string    `json:"-"`
	ToAccountId       *int64     `json:"toAccountId"`
	ToBankId          *int64     `json:"-"`
	ToAccountName     *string    `json:"-"`
	ToAccountNumber   *string    `json:"-"`
	CreditAmount      float32    `json:"creditAmount" validate:"required"`
	PaidAmount        float32    `json:"-"`
	DepositChannel    string     `json:"depositChannel"`
	OverAmount        float32    `json:"overAmount"`
	BonusAmount       float32    `json:"bonusAmount"`
	BeforeAmount      float32    `json:"-"`
	AfterAmount       float32    `json:"-"`
	TransferAt        *time.Time `json:"transferAt" example:"2023-05-31T22:33:44+07:00"`
	CreatedByUserId   int64      `json:"-"`
	CreatedByUsername string     `json:"-"`
	Status            string     `json:"-"`
	IsAutoCredit      bool       `json:"isAutoCredit"`
}

type BonusTransactionCreateBody struct {
	MemberCode        string    `json:"memberCode" validate:"required"`
	UserId            int64     `json:"-"`
	TransferType      string    `json:"-"`
	ToAccountId       int64     `json:"-"`
	ToBankId          int64     `json:"-"`
	ToAccountName     string    `json:"-"`
	ToAccountNumber   string    `json:"-"`
	BonusAmount       float32   `json:"bonusAmount" validate:"required"`
	BonusReason       string    `json:"bonusReason"`
	BeforeAmount      float32   `json:"-"`
	AfterAmount       float32   `json:"-"`
	TransferAt        time.Time `json:"transferAt" validate:"required" example:"2023-05-31T22:33:44+07:00"`
	CreatedByUserId   int64     `json:"-"`
	CreatedByUsername string    `json:"-"`
	Status            string    `json:"-"`
}

type BankTransactionUpdateBody struct {
	Status            string    `json:"-" validate:"required"`
	RemovedAt         time.Time `json:"removedAt" example:"2023-05-31T22:33:44+07:00"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}

type BankTransactionResponse struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	UserId              int64          `json:"userId"`
	MemberCode          string         `json:"memberCode"`
	UserUsername        string         `json:"userUsername"`
	UserFullname        string         `json:"userFullname"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float32        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float32        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float32        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float32        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float32        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float32        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float32        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          time.Time      `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         time.Time      `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           time.Time      `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}

type PendingDepositTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type PendingWithdrawTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type FinishedTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:1"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:2"`
	AccountId        string `form:"accountId" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type RemovedTransactionListRequest struct {
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:1"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:2"`
	AccountId        string `form:"accountId" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type BankTransactionCancelBody struct {
	Status             string    `json:"-"`
	CancelRemark       string    `json:"cancelRemark" validate:"required"`
	CanceledAt         time.Time `json:"-"`
	CanceledByUserId   int64     `json:"-"`
	CanceledByUsername string    `json:"-"`
}

type BankConfirmDepositRequest struct {
	TransferAt          *time.Time `json:"transferAt" validate:"required"`
	SlipUrl             string     `json:"slipUrl" validate:"required"`
	BonusAmount         float32    `json:"bonusAmount" validate:"required"`
	ConfirmedAt         time.Time  `json:"-"`
	ConfirmedByUserId   int64      `json:"-"`
	ConfirmedByUsername string     `json:"-"`
}

type BankConfirmWithdrawRequest struct {
	CreditAmount        float32   `json:"creditAmount" validate:"required"`
	BankChargeAmount    float32   `json:"bankChargeAmount" validate:"required"`
	ConfirmedAt         time.Time `json:"-"`
	ConfirmedByUserId   int64     `json:"-"`
	ConfirmedByUsername string    `json:"-"`
}

type BankDepositTransactionConfirmBody struct {
	TransferAt          time.Time `json:"transferAt"`
	BonusAmount         float32   `json:"bonusAmount"`
	Status              string    `json:"status"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type BankWithdrawTransactionConfirmBody struct {
	TransferAt          time.Time `json:"transferAt"`
	CreditAmount        float32   `json:"creditAmount"`
	BankChargeAmount    float32   `json:"bankChargeAmount"`
	Status              string    `json:"status"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type CreateBankTransactionActionBody struct {
	TransactionId       int64     `json:"transactionId"`
	UserId              int64     `json:"userId"`
	TransferType        string    `json:"transferType"`
	FromAccountId       int64     `json:"fromAccountId"`
	ToAccountId         int64     `json:"toAccountId"`
	JsonBefore          string    `json:"jsonBefore"`
	TransferAt          time.Time `json:"transferAt"`
	SlipUrl             string    `json:"slipUrl"`
	BonusAmount         float32   `json:"bonusAmount"`
	CreditAmount        float32   `json:"creditAmount"`
	BankChargeAmount    float32   `json:"bankChargeAmount"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type CreateBankStatementActionBody struct {
	StatementId         int64     `json:"statementId"`
	UserId              int64     `json:"userId"`
	ActionType          string    `json:"actionType"`
	AccountId           int64     `json:"accountId"`
	JsonBefore          string    `json:"jsonBefore"`
	ConfirmedAt         time.Time `json:"confirmedAt"`
	ConfirmedByUserId   int64     `json:"confirmedByUserId"`
	ConfirmedByUsername string    `json:"confirmedByUsername"`
}

type BankTransactionRemoveBody struct {
	Status            string    `json:"-" validate:"required"`
	RemovedAt         time.Time `json:"removedAt"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}

type MemberTransactionListRequest struct {
	UserId           string `form:"userId" extensions:"x-order:1"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:2"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:3"`
	TransferType     string `form:"transferType" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type MemberTransactionSummary struct {
	TotalDepositAmount  float32 `json:"totalDepositAmount"`
	TotalWithdrawAmount float32 `json:"totalWithdrawAmount"`
	TotalBonusAmount    float32 `json:"totalBonusAmount"`
}

type MemberTransaction struct {
	Id                  int64          `json:"id" gorm:"primaryKey"`
	UserId              int64          `json:"userId"`
	MemberCode          string         `json:"memberCode"`
	UserUsername        string         `json:"userUsername"`
	UserFullname        string         `json:"userFullname"`
	TransferType        string         `json:"transferType"`
	PromotionId         int64          `json:"promotionId"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FromBankName        string         `json:"fromBankName"`
	FromAccountName     string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToBankName          string         `json:"toBankName"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	CreditAmount        float32        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount          float32        `json:"paidAmount" sql:"type:decimal(14,2);"`
	DepositChannel      string         `json:"depositChannel"`
	OverAmount          float32        `json:"overAmount" sql:"type:decimal(14,2);"`
	BonusAmount         float32        `json:"bonusAmount" sql:"type:decimal(14,2);"`
	BonusReason         string         `json:"bonusReason"`
	BeforeAmount        float32        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount         float32        `json:"afterAmount" sql:"type:decimal(14,2);"`
	BankChargeAmount    float32        `json:"bankChargeAmount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreatedByUserId     int64          `json:"createdByUserId"`
	CreatedByUsername   string         `json:"createdByUsername"`
	CancelRemark        string         `json:"cancelRemark"`
	CanceledAt          time.Time      `json:"canceledAt"`
	CanceledByUserId    int64          `json:"canceledByUserId"`
	CanceledByUsername  string         `json:"canceledByUsername"`
	ConfirmedAt         *time.Time     `json:"confirmedAt"`
	ConfirmedByUserId   int64          `json:"confirmedByUserId"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	RemovedAt           time.Time      `json:"removedAt"`
	RemovedByUserId     int64          `json:"removedByUserId"`
	RemovedByUsername   string         `json:"removedByUsername"`
	Status              string         `json:"status"`
	StatusDetail        string         `json:"statusDetail"`
	IsAutoCredit        bool           `json:"isAutoCredit"`
	CreatedAt           time.Time      `json:"createAt"`
	UpdatedAt           *time.Time     `json:"updateAt"`
	DeletedAt           gorm.DeletedAt `json:"deleteAt"`
}

type BankAutoDepositCondition struct {
	Id              int64   `json:"-"`
	UserId          int64   `json:"-"`
	ToAccountId     int64   `json:"toAccountId"`
	MinCreditAmount float32 `json:"minCreditAmount"`
	MaxCreditAmount float32 `json:"maxCreditAmount"`
}

type BankAutoWithdrawCondition struct {
	Id                      int64   `json:"-"`
	UserId                  int64   `json:"-"`
	FromAccountId           int64   `json:"toAccountId"`
	MinCreditAmount         float32 `json:"minCreditAmount"`
	MaxCreditAmount         float32 `json:"maxCreditAmount"`
	AutoWithdrawCreditFlag  string  `json:"autoWithdrawCreditFlag"`
	AutoWithdrawConfirmFlag string  `json:"autoWithdrawConfirmFlag"`
}
