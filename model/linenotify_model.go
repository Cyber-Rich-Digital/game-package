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

type LinenotifyGame struct {
	Id           int64      `json:"id"`
	Name         string     `json:"name" validate:"required"`
	ClientId     string     `json:"clientid" validate:"required"`
	ClientSecret string     `json:"clientsecret" validate:"required"`
	ResponseType string     `json:"responsetype" validate:"required"`
	RedirectUri  string     `json:"redirecturi" validate:"required"`
	Scope        string     `json:"scope" validate:"required"`
	State        string     `json:"state" validate:"required"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type LinenotifyGameResponse struct {
	Id           int64  `json:"id"`
	Name         string `json:"name" validate:"required"`
	Clientid     string `json:"clientid" validate:"required"`
	Clientsecret string `json:"clientsecret" validate:"required"`
	Responsetype string `json:"responsetype" validate:"required"`
	Redirecturi  string `json:"redirecturi" validate:"required"`
	Scope        string `json:"scope" validate:"required"`
	State        string `json:"state" validate:"required"`
}

type LinenotifyGameParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LineNoifyUsergame struct {
	UserId       int64      `json:"name" validate:"required"`
	TypeNotifyId string     `json:"TypeNotifyId" validate:"required"`
	Token        string     `json:"token" validate:"required"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

type LineNoifyUsergameBody struct {
	UserId       int64  `json:"name" validate:"required"`
	TypeNotifyId string `json:"TypeNotifyId" validate:"required"`
	Token        string `json:"token" validate:"required"`
}

type LineNotifyUserGameParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type LineNotifyUserGametDeleteBody struct {
	Id        string    `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type LineNotifyRedirectReponse struct {
	Code  string `json:"code"`
	State string `json:"state"`
}
