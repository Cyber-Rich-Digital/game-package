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

type BankAccount struct {
	Id                    int64          `json:"id"`
	BankId                int64          `json:"bankId"`
	AccountTypeId         int64          `json:"accountTypeId"`
	AccountName           string         `json:"accountHame"`
	AccountNumber         string         `json:"accountNumber"`
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
	UpdatedAt             time.Time      `json:"updatedAt"`
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
	AccountName   string `json:"accountName" validate:"required"`
	AccountNumber string `json:"domainName" validate:"required"`
}

type BankAccountResponse struct {
	Id               int64          `json:"id"`
	BankId           int64          `json:"bankId"`
	AccountTypeId    int64          `json:"accountTypeId"`
	AccountName      string         `json:"accountHame"`
	AccountNumber    string         `json:"accountNumber"`
	TransferPriority string         `json:"transferPriority"`
	AccountStatus    string         `json:"accountStatus"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"deletedAt"`
}
