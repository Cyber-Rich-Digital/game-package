package model

import (
	"time"

	"gorm.io/gorm"
)

type Website struct {
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

type WebsiteParam struct {
	Id int `uri:"id" binding:"required"`
}

type WebsiteDate struct {
	Date   int64 `uri:"date" binding:"required"`
	UserId int
}

type WebsiteQuery struct {
	Page   int    `form:"page" default:"1" min:"1"`
	Limit  int    `form:"limit" default:"10" min:"1" max:"100"`
	Search string `form:"search"`
	Sort   int    `form:"sort"`
	UserId int
	Role   string
}

type WebsiteBody struct {
	Title      string `json:"title" validate:"required"`
	DomainName string `json:"domainName" validate:"required"`
}

type WebsiteResponse struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	DomainName string    `json:"domainName"`
	ApiKey     string    `json:"apikey"`
	CreatedAt  time.Time `json:"createdAt"`
	Total      int       `json:"total"`
}

type WebsiteListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type WebsiteList struct {
	Id         int    `json:"id" gorm:"column:id"`
	Title      string `json:"title" gorm:"column:title"`
	DomainName string `json:"domainName" gorm:"column:domain_name"`
	UserId     uint   `json:"userId" gorm:"column:user_id"`
}
