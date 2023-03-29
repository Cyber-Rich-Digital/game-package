package model

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	IconUrl   string    `json:"iconUrl"`
	TypeFlag  string    `json:"typeFlag"`
	CreatedAt time.Time `json:"createdAt"`
}
type BankResponse struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code"`
	IconUrl  string `json:"iconUrl"`
	TypeFlag string `json:"typeFlag"`
}
type BankListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}
type BankListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type AccountType struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	TypeFlag  string    `json:"typeFlag"`
	CreatedAt time.Time `json:"createdAt"`
}
type AccountTypeResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
type AccountTypeListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}
type AccountTypeListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type BankAccount struct {
	Id                    int64          `json:"id"`
	BankId                int64          `json:"bankId"`
	BankName              string         `json:"bankName"`
	AccountTypeId         int64          `json:"accountTypeId"`
	AccountTypeName       string         `json:"accountTypeName"`
	AccountName           string         `json:"accountHame"`
	AccountNumber         string         `json:"accountNumber"`
	AccountBalance        float32        `json:"accountBalance" sql:"type:decimal(14,2);"`
	TransferPriority      string         `json:"transferPriority"`
	AccountStatus         string         `json:"accountStatus"`
	DeviceUid             string         `json:"deviceUid"`
	PinCode               string         `json:"pinCode"`
	ConectionStatus       string         `json:"conectionStatus"`
	AutoCreditFlag        string         `json:"autoCreditFlag"`
	AutoWithdrawFlag      string         `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string         `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string         `json:"autoTransferMaxAmount"`
	QrWalletStatus        string         `json:"qrWalletStatus"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             *time.Time     `json:"updatedAt"`
	DeletedAt             gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type BankAccountBody struct {
	BankId                int64  `json:"bankId" validate:"required"`
	AccountTypeId         int64  `json:"AccounTypeId" validate:"required"`
	AccountName           string `json:"accountName" validate:"required"`
	AccountNumber         string `json:"accountNumber" validate:"required"`
	DeviceUid             string `json:"deviceUid"`
	PinCode               string `json:"pinCode"`
	AutoCreditFlag        string `json:"autoCreditFlag"`
	AutoWithdrawFlag      string `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string `json:"autoTransferMaxAmount"`
	TransferPriority      string `json:"transferPriority"`
	QrWalletStatus        string `json:"qrWalletStatus"`
}

type BankAccountResponse struct {
	Id               int64          `json:"id"`
	BankId           int64          `json:"bankId"`
	BankName         string         `json:"bankName"`
	AccountTypeId    int64          `json:"accountTypeId"`
	AccountTypeName  string         `json:"accountTypeName"`
	AccountName      string         `json:"accountHame"`
	AccountNumber    string         `json:"accountNumber"`
	AccountBalance   float32        `json:"accountBalance"`
	TransferPriority string         `json:"transferPriority"`
	AccountStatus    string         `json:"accountStatus"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        *time.Time     `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransaction struct {
	Id               int64          `json:"id"`
	AccountId        int64          `json:"accountId"`
	Description      string         `json:"description"`
	TransferType     string         `json:"transferType"`
	Amount           float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt       time.Time      `json:"transferAt"`
	CreateByUsername string         `json:"createByUsername"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        *time.Time     `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransactionParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransactionListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type BankAccountTransactionBody struct {
	AccountId    int64     `json:"accountId" validate:"required"`
	Description  string    `json:"description"`
	TransferType string    `json:"transferType" validate:"required"`
	Amount       float32   `json:"amount" validate:"required"`
	TransferAt   time.Time `json:"transferAt" validate:"required"`
}

type BankAccountTransactionResponse struct {
	Id               int64          `json:"id"`
	AccountId        int64          `json:"accountId"`
	Description      string         `json:"description"`
	TransferType     string         `json:"transferType"`
	Amount           float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt       time.Time      `json:"transferAt"`
	CreateByUsername string         `json:"createByUsername"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        *time.Time     `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransfer struct {
	Id                  int64          `json:"id"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FfromAccountName    string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	Amount              float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreateByUsername    string         `json:"createByUsername"`
	Status              string         `json:"status"`
	ConfirmedAt         time.Time      `json:"confirmedAt"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           *time.Time     `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransferParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransferListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type BankAccountTransferBody struct {
	Status        string    `json:"status" validate:"required"`
	FromAccountId int64     `json:"fromAccountId" validate:"required"`
	ToAccountId   int64     `json:"toAccountId" validate:"required"`
	Amount        float32   `json:"amount" validate:"required"`
	TransferAt    time.Time `json:"transferAt" validate:"required"`
}

type BankAccountTransferConfirmBody struct {
	Status              string    `json:"status" validate:"required"`
	ConfirmedByUsername string    `json:"confirmedByUsername" validate:"required"`
	ConfirmedAt         time.Time `json:"confirmedAt" validate:"required"`
}

type BankAccountTransferResponse struct {
	Id                  int64          `json:"id"`
	FromAccountId       int64          `json:"fromAccountId"`
	FromBankId          int64          `json:"fromBankId"`
	FfromAccountName    string         `json:"fromAccountName"`
	FromAccountNumber   string         `json:"fromAccountNumber"`
	ToAccountId         int64          `json:"toAccountId"`
	ToBankId            int64          `json:"toBankId"`
	ToAccountName       string         `json:"toAccountName"`
	ToAccountNumber     string         `json:"toAccountNumber"`
	Amount              float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt          time.Time      `json:"transferAt"`
	CreateByUsername    string         `json:"createByUsername"`
	Status              string         `json:"status"`
	ConfirmedAt         time.Time      `json:"confirmedAt"`
	ConfirmedByUsername string         `json:"confirmedByUsername"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           *time.Time     `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `json:"deletedAt"`
}
