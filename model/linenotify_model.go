package model

import (
	"time"
)

type Linenotify struct {
	Id          int64      `json:"id"`
	StartCredit float32    `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string     `json:"token" validate:"required"`
	NotifyId    int64      `json:"notifyId" validate:"required"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}
type LinenotifyResponse struct {
	Id          int64   `json:"id"`
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}
type LinenotifyListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type LinenotifyParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LinenotifyListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type LinenotifyCreateBody struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}
type LinenotifyUpdateBody struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}

type LinenotifyUpdateRequest struct {
	StartCredit float32 `json:"startcredit" sql:"type:decimal(14,2);"`
	Token       string  `json:"token" validate:"required"`
	NotifyId    int64   `json:"notifyId" validate:"required"`
	Status      string  `json:"status"`
}
