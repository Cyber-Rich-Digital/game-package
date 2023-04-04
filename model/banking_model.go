package model

import (
	"time"

	"gorm.io/gorm"
)

type BankStatement struct {
	Id         int64          `json:"id" gorm:"primaryKey"`
	AccountId  int64          `json:"accountId"`
	Amount     float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt time.Time      `json:"transferAt"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"createAt"`
	UpdatedAt  *time.Time     `json:"updateAt"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt"`
}

type BankStatementGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankStatementListRequest struct {
	AccountId        string `form:"accountId" extensions:"x-order:1"`
	Amount           string `form:"amount" extensions:"x-order:2"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type BankStatementCreateBody struct {
	AccountId  int64     `json:"accountId"`
	Amount     float32   `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt time.Time `json:"transferAt"`
	Status     string    `json:"-"`
}

type BankStatementUpdateBody struct {
	Status string `json:"status" validate:"required"`
}

type BankStatementResponse struct {
	Id         int64          `json:"id" gorm:"primaryKey"`
	AccountId  int64          `json:"accountId"`
	Amount     float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt time.Time      `json:"transferAt"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"createAt"`
	UpdatedAt  *time.Time     `json:"updateAt"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt"`
}

type BankTransaction struct {
	Id                int64          `json:"id" gorm:"primaryKey"`
	MemberCode        string         `json:"memberCode"`
	UserId            int64          `json:"userId"`
	TransferType      string         `json:"transferType"`
	PromotionId       int64          `json:"promotionId"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromAccountName   string         `json:"fromAccountName"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	CreditAmount      float32        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount        float32        `json:"paidAmount" sql:"type:decimal(14,2);"`
	OverAmount        float32        `json:"overAmount" sql:"type:decimal(14,2);"`
	BeforeAmount      float32        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount       float32        `json:"afterAmount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUserId   int64          `json:"createdByUserId"`
	CreatedByUsername string         `json:"createdByUsername"`
	RemovedAt         time.Time      `json:"removedAt"`
	RemovedByUserId   int64          `json:"removedByUserId"`
	RemovedByUsername string         `json:"removedByUsername"`
	Status            string         `json:"status"`
	IsAutoCredit      bool           `json:"isAutoCredit"`
	CreatedAt         time.Time      `json:"createAt"`
	UpdatedAt         *time.Time     `json:"updateAt"`
	DeletedAt         gorm.DeletedAt `json:"deleteAt"`
}

type BankTransactionGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankTransactionListRequest struct {
	MemberCode       string `form:"memberCode" extensions:"x-order:1"`
	UserId           string `form:"userId" extensions:"x-order:2"`
	FromTransferDate string `form:"fromTransferDate" extensions:"x-order:3"`
	ToTransferDate   string `form:"toTransferDate" extensions:"x-order:4"`
	Search           string `form:"search" extensions:"x-order:5"`
	Page             int    `form:"page" extensions:"x-order:6" default:"1" min:"1"`
	Limit            int    `form:"limit" extensions:"x-order:7" default:"10" min:"1" max:"100"`
	SortCol          string `form:"sortCol" extensions:"x-order:8"`
	SortAsc          string `form:"sortAsc" extensions:"x-order:9"`
}

type BankTransactionCreateBody struct {
	MemberCode        string    `json:"memberCode" validate:"required"`
	UserId            int64     `json:"-"`
	TransferType      string    `json:"transferType" validate:"required"`
	PromotionId       int64     `json:"promotionId"`
	FromAccountId     int64     `json:"fromAccountId"`
	FromBankId        int64     `json:"-"`
	FromAccountName   string    `json:"-"`
	FromAccountNumber string    `json:"-"`
	ToAccountId       int64     `json:"toAccountId"`
	ToBankId          int64     `json:"-"`
	ToAccountName     string    `json:"-"`
	ToAccountNumber   string    `json:"-"`
	CreditAmount      float32   `json:"creditAmount"`
	PaidAmount        float32   `json:"paidAmount"`
	OverAmount        float32   `json:"overAmount"`
	BeforeAmount      float32   `json:"beforeAmount"`
	AfterAmount       float32   `json:"afterAmount"`
	TransferAt        time.Time `json:"transferAt"`
	CreatedByUserId   int64     `json:"-"`
	CreatedByUsername string    `json:"-"`
	Status            string    `json:"-"`
	IsAutoCredit      bool      `json:"isAutoCredit" validate:"required"`
}

type BankTransactionUpdateBody struct {
	Status            string    `json:"-" validate:"required"`
	RemovedAt         time.Time `json:"removedAt"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}

type BankTransactionResponse struct {
	Id                int64          `json:"id" gorm:"primaryKey"`
	MemberCode        string         `json:"memberCode"`
	UserId            int64          `json:"userId"`
	TransferType      string         `json:"transferType"`
	PromotionId       int64          `json:"promotionId"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromAccountName   string         `json:"fromAccountName"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	CreditAmount      float32        `json:"creditAmount" sql:"type:decimal(14,2);"`
	PaidAmount        float32        `json:"paidAmount" sql:"type:decimal(14,2);"`
	OverAmount        float32        `json:"overAmount" sql:"type:decimal(14,2);"`
	BeforeAmount      float32        `json:"beforeAmount" sql:"type:decimal(14,2);"`
	AfterAmount       float32        `json:"afterAmount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUserId   int64          `json:"createdByUserId"`
	CreatedByUsername string         `json:"createdByUsername"`
	RemovedAt         time.Time      `json:"removedAt"`
	RemovedByUserId   int64          `json:"removedByUserId"`
	RemovedByUsername string         `json:"removedByUsername"`
	Status            string         `json:"status"`
	CreatedAt         time.Time      `json:"createAt"`
	UpdatedAt         *time.Time     `json:"updateAt"`
	DeletedAt         gorm.DeletedAt `json:"deleteAt"`
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

type BankTransactionRemoveBody struct {
	Status            string    `json:"-" validate:"required"`
	RemovedAt         time.Time `json:"removedAt"`
	RemovedByUserId   int64     `json:"removedByUserId"`
	RemovedByUsername string    `json:"removedByUsername"`
}
