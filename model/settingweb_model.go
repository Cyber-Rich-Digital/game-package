package model

import (
	"time"
)

type Settingweb struct {
	Id             int64     `json:"id"`
	Logo           string    `json:"logo"`
	BackgrondColor string    `json:"backgrondcolor"`
	UserAuto       string    `json:"userAuto"`
	OtpRegister    string    `json:"otpRegister"`
	TranWithdraw   string    `json:"tranWithdraw"`
	Register       string    `json:"register"`
	DepositFirst   string    `json:"depositFirst"`
	DepositNext    string    `json:"depositNext"`
	Withdraw       string    `json:"withdraw"`
	Line           string    `json:"line"`
	Url            string    `json:"url"`
	Opt            string    `json:"opt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	DeletedAt      time.Time `json:"deletedAt"`
}
type SettingwebResponse struct {
	Id             int64  `json:"id"`
	Logo           string `json:"logo"`
	BackgrondColor string `json:"backgrondcolor"`
	UserAuto       string `json:"userAuto"`
	OtpRegister    string `json:"otpRegister"`
	TranWithdraw   string `json:"tranWithdraw"`
	Register       string `json:"register"`
	DepositFirst   string `json:"depositFirst"`
	DepositNext    string `json:"depositNext"`
	Withdraw       string `json:"withdraw"`
	Line           string `json:"line"`
	Url            string `json:"url"`
	Opt            string `json:"opt"`
}
type SettingwebListResponse struct {
	Id    int `json:"id"`
	Total int `json:"total"`
}

type SettingwebParam struct {
	Id int64 `uri:"id" binding:"required"`
}

type SettingwebListRequest struct {
	Page    int    `form:"page" default:"1" min:"1"`
	Limit   int    `form:"limit" default:"10" min:"1" max:"100"`
	Search  string `form:"search"`
	SortCol string `form:"sortCol"`
	SortAsc string `form:"sortAsc"`
}

type SettingwebCreateBody struct {
	Id             int64  `json:"id"`
	Logo           string `json:"logo"`
	BackgrondColor string `json:"backgrondcolor"`
	UserAuto       string `json:"userAuto"`
	OtpRegister    string `json:"otpRegister"`
	TranWithdraw   string `json:"tranWithdraw"`
	Register       string `json:"register"`
	DepositFirst   string `json:"depositFirst"`
	DepositNext    string `json:"depositNext"`
	Withdraw       string `json:"withdraw"`
	Line           string `json:"line"`
	Url            string `json:"url"`
	Opt            string `json:"opt"`
}
type SettingwebUpdateBody struct {
	Id             int64  `json:"id"`
	Logo           string `json:"logo"`
	BackgrondColor string `json:"backgrondcolor"`
	UserAuto       string `json:"userAuto"`
	OtpRegister    string `json:"otpRegister"`
	TranWithdraw   string `json:"tranWithdraw"`
	Register       string `json:"register"`
	DepositFirst   string `json:"depositFirst"`
	DepositNext    string `json:"depositNext"`
	Withdraw       string `json:"withdraw"`
	Line           string `json:"line"`
	Url            string `json:"url"`
	Opt            string `json:"opt"`
}
