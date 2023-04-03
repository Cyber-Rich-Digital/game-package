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
	UpdatedAt  time.Time      `json:"updateAt"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt"`
}

type BankStatementGetRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankStatementListRequest struct {
	AccountId        int     `form:"accountId"`
	Amount           float32 `form:"amount"`
	FromTransferDate string  `form:"fromTransferDate"`
	ToTransferDate   string  `form:"toTransferDate"`
	Search           string  `form:"search"`
	Page             int     `form:"page" default:"1" min:"1"`
	Limit            int     `form:"limit" default:"10" min:"1" max:"100"`
	SortCol          string  `form:"sortCol"`
	SortAsc          string  `form:"sortAsc"`
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
	UpdatedAt  time.Time      `json:"updateAt"`
	DeletedAt  gorm.DeletedAt `json:"deleteAt"`
}
