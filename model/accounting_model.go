package model

import (
	"time"

	"gorm.io/gorm"
)

type SimpleOption struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type ConfirmRequest struct {
	Password string `json:"password"`
	UserId   int64  `json:"-"`
}

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
	BankIconUrl           string         `json:"bankIconUrl"`
	AccountTypeId         int64          `json:"accountTypeId"`
	AccountTypeName       string         `json:"accountTypeName"`
	AccountName           string         `json:"accountHame"`
	AccountNumber         string         `json:"accountNumber"`
	AccountBalance        float32        `json:"accountBalance" sql:"type:decimal(14,2);"`
	AccountPriority       string         `json:"accountPriority"`
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

type BankAccountCreateBody struct {
	BankId                int64   `json:"bankId" validate:"required"`
	AccountTypeId         int64   `json:"accounTypeId" validate:"required"`
	AccountName           string  `json:"accountName" validate:"required"`
	AccountNumber         string  `json:"accountNumber" validate:"required"`
	AccountBalance        float32 `json:"-"`
	DeviceUid             string  `json:"deviceUid"`
	PinCode               string  `json:"pinCode"`
	AutoCreditFlag        string  `json:"autoCreditFlag"`
	AutoWithdrawFlag      string  `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string  `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string  `json:"autoTransferMaxAmount"`
	AccountPriority       string  `json:"accountPriority"`
	QrWalletStatus        string  `json:"qrWalletStatus"`
	AccountStatus         string  `json:"accountStatus"`
	ConectionStatus       string  `json:"-"`
}
type BankAccountUpdateBody struct {
	BankId                int64  `json:"-"`
	AccountTypeId         int64  `json:"accounTypeId" validate:"required"`
	AccountName           string `json:"-"`
	AccountNumber         string `json:"-"`
	DeviceUid             string `json:"deviceUid"`
	PinCode               string `json:"pinCode"`
	AutoCreditFlag        string `json:"autoCreditFlag"`
	AutoWithdrawFlag      string `json:"autoWithdrawFlag"`
	AutoWithdrawMaxAmount string `json:"autoWithdrawMaxAmount"`
	AutoTransferMaxAmount string `json:"autoTransferMaxAmount"`
	AccountPriority       string `json:"accountPriority"`
	QrWalletStatus        string `json:"qrWalletStatus"`
	AccountStatus         string `json:"accountStatus"`
	ConectionStatus       string `json:"-"`
}

type BankAccountResponse struct {
	Id              int64          `json:"id"`
	BankId          int64          `json:"bankId"`
	BankName        string         `json:"bankName"`
	BankIconUrl     string         `json:"bankIconUrl"`
	AccountTypeId   int64          `json:"accountTypeId"`
	AccountTypeName string         `json:"accountTypeName"`
	AccountName     string         `json:"accountHame"`
	AccountNumber   string         `json:"accountNumber"`
	AccountBalance  float32        `json:"accountBalance"`
	AccountPriority string         `json:"accountPriority"`
	AccountStatus   string         `json:"accountStatus"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       *time.Time     `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransaction struct {
	Id                int64          `json:"id"`
	AccountId         int64          `json:"accountId"`
	Description       string         `json:"description"`
	TransferType      string         `json:"transferType"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransactionParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransactionListRequest struct {
	AccountId       int    `form:"accountId"`
	FromCreatedDate string `form:"fromCreatedDate"`
	TocreatedDate   string `form:"toCreatedDate"`
	TransferType    string `form:"transferType"`
	Search          string `form:"search"`
	Page            int    `form:"page" default:"1" min:"1"`
	Limit           int    `form:"limit" default:"10" min:"1" max:"100"`
	SortCol         string `form:"sortCol"`
	SortAsc         string `form:"sortAsc"`
}

type BankAccountTransactionBody struct {
	AccountId         int64     `json:"accountId" validate:"required"`
	Description       string    `json:"description"`
	TransferType      string    `json:"transferType" validate:"required"`
	Amount            float32   `json:"amount" validate:"required"`
	TransferAt        time.Time `json:"transferAt" validate:"required"`
	CreatedByUsername string    `json:"-"`
}

type BankAccountTransactionResponse struct {
	Id                int64          `json:"id"`
	AccountId         int64          `json:"accountId"`
	BankName          string         `json:"bankName"`
	AccountName       string         `json:"accountHame"`
	AccountNumber     string         `json:"accountNumber"`
	Description       string         `json:"description"`
	TransferType      string         `json:"transferType"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransfer struct {
	Id                int64          `json:"id"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankName      string         `json:"fromBankName"`
	FromAccountName   string         `json:"fromAccountName"`
	FromAccountNumber string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	TobankName        string         `json:"toBankName"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	Status            string         `json:"status"`
	ConfirmedAt       time.Time      `json:"confirmedAt"`
	ConfirmedByUserId int64          `json:"confirmedByUserId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountTransferParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type BankAccountTransferListRequest struct {
	AccountId       int    `form:"accountId"`
	FromCreatedDate string `form:"fromCreatedDate"`
	TocreatedDate   string `form:"toCreatedDate"`
	ToAccountId     int    `form:"toAccountId"`
	Search          string `form:"search"`
	Page            int    `form:"page" default:"1" min:"1"`
	Limit           int    `form:"limit" default:"10" min:"1" max:"100"`
	SortCol         string `form:"sortCol"`
	SortAsc         string `form:"sortAsc"`
}

type BankAccountTransferBody struct {
	Status            string    `json:"-"`
	FromAccountId     int64     `json:"fromAccountId" validate:"required"`
	FromBankId        int64     `json:"-"`
	FromAccountName   string    `json:"-"`
	FromAccountNumber string    `json:"-"`
	ToAccountId       int64     `json:"toAccountId" validate:"required"`
	ToBankId          int64     `json:"-"`
	ToAccountName     string    `json:"-"`
	ToAccountNumber   string    `json:"-"`
	Amount            float32   `json:"amount" validate:"required"`
	TransferAt        time.Time `json:"transferAt" validate:"required"`
	CreatedByUsername string    `json:"-"`
}

type BankAccountTransferConfirmBody struct {
	Status            string    `json:"status" validate:"required"`
	ConfirmedByUserId int64     `json:"confirmedByUserId" validate:"required"`
	ConfirmedAt       time.Time `json:"confirmedAt" validate:"required"`
}

type BankAccountTransferResponse struct {
	Id                int64          `json:"id"`
	FromAccountId     int64          `json:"fromAccountId"`
	FromBankId        int64          `json:"fromBankId"`
	FromBankName      string         `json:"fromBankName"`
	FromAccountName   string         `json:"fromAccountName"`
	FomAccountNumber  string         `json:"fromAccountNumber"`
	ToAccountId       int64          `json:"toAccountId"`
	ToBankId          int64          `json:"toBankId"`
	ToBankName        string         `json:"toBankName"`
	ToAccountName     string         `json:"toAccountName"`
	ToAccountNumber   string         `json:"toAccountNumber"`
	Amount            float32        `json:"amount" sql:"type:decimal(14,2);"`
	TransferAt        time.Time      `json:"transferAt"`
	CreatedByUsername string         `json:"createdByUsername"`
	Status            string         `json:"status"`
	ConfirmedAt       time.Time      `json:"confirmedAt"`
	ConfirmedByUserId int64          `json:"confirmedByUserId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         *time.Time     `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt"`
}