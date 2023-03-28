package model

import (
	"time"

	"gorm.io/gorm"
)

type BankAccount struct {
	Id         int            `json:"id"`
	Title      string         `json:"title"`
	DomainName string         `json:"domainName"`
	ApiKey     string         `json:"apikey"`
	UserId     int            `json:"userId"`
	Tags       []TagList      `json:"tags"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt"`
}

type BankAccountParam struct {
	Id int `uri:"id" binding:"required"`
}

type BankAccountDate struct {
	Date   int64 `uri:"date" binding:"required"`
	UserId int
}

type BankAccountQuery struct {
	Page   int    `form:"page" default:"1" min:"1"`
	Limit  int    `form:"limit" default:"10" min:"1" max:"100"`
	Search string `form:"search"`
	Sort   int    `form:"sort"`
	UserId int
	Role   string
}

type BankAccountBody struct {
	Title      string `json:"title" validate:"required"`
	DomainName string `json:"domainName" validate:"required"`
}

type BankAccountResponse struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	DomainName string    `json:"domainName"`
	ApiKey     string    `json:"apikey"`
	CreatedAt  time.Time `json:"createdAt"`
	Total      int       `json:"total"`
}

type BankAccountListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type BankAccountList struct {
	Id         int    `json:"id" gorm:"column:id"`
	Title      string `json:"title" gorm:"column:title"`
	DomainName string `json:"domainName" gorm:"column:domain_name"`
	UserId     uint   `json:"userId" gorm:"column:user_id"`
}
