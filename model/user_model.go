package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id            int64          `json:"id"`
	Partner       string         `json:"partner"`
	MemberCode    string         `json:"memberCode"`
	Username      string         `json:"username"`
	Phone         string         `json:"phone"`
	Promotion     string         `json:"promotion"`
	Password      string         `json:"password"`
	Status        string         `json:"status"`
	Firstname     string         `json:"firstname"`
	Lastname      string         `json:"lastname"`
	Fullname      string         `json:"fullname"`
	Bankname      string         `json:"bankname"`
	BankAccount   string         `json:"bankAccount"`
	Channel       string         `json:"channel"`
	TrueWallet    string         `json:"trueWallet"`
	Contact       string         `json:"contact"`
	Note          string         `json:"note"`
	Course        string         `json:"course"`
	Credit        int            `json:"credit"`
	TurnoverLimit int            `json:"turnoverLimit"`
	Ip            string         `json:"ip"`
	IpRegistered  string         `json:"ipRegistered"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"`
	LogedinAt     *time.Time     `json:"logedinAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type CreateUser struct {
	Partner      string `json:"partner" validate:"required,max=20"`
	MemberCode   string `json:"memberCode" validate:"max=255"`
	Username     string `json:"username" validate:"required,max=255"`
	Phone        string `json:"phone" validate:"required,max=12" example:"0812345678"`
	Promotion    string `json:"promotion" validate:"max=20"`
	Password     string `json:"password" validate:"required,max=255"`
	Status       string `json:"status" validate:"max=255" enum:"ACTIVE,DEACTIVE" example:"ACTIVE"`
	Fullname     string `json:"fullname" validate:"required,max=255"`
	Bankname     string `json:"bankname" validate:"required,max=50"`
	BankAccount  string `json:"bankAccount" validate:"required,max=15"`
	Channel      string `json:"channel" validate:"required,max=20" enum:"Google,Youtube,Facebook" example:"Google"`
	TrueWallet   string `json:"trueWallet" validate:"required,max=20"`
	Contact      string `json:"contact" validate:"max=255"`
	Note         string `json:"note" validate:"max=255"`
	Course       string `json:"course" validate:"max=50"`
	IpRegistered string `json:"ipRegistered" validate:"required,max=20" example:"1.1.1.1"`
}

type LoginUser struct {
	Username string `json:"username" validate:"required,min=8,max=30"`
	Password string `json:"password" validate:"required,8,max=30"`
	IP       string `json:"ip"`
}

type UserLoginUpdate struct {
	IP        string    `json:"ip"`
	LogedinAt time.Time `json:"logedinAt"`
}

type UserListQuery struct {
	Page   int    `form:"page" validate:"required,min=1"`
	Limit  int    `form:"limit" validate:"required,min=1,max=100"`
	Search string `form:"search"`
	Status string `form:"status"`
}

type UpdateUser struct {
	Partner     string `json:"partner" validate:"required,max=20"`
	MemberCode  string `json:"memberCode" validate:"max=255"`
	Promotion   string `json:"promotion" validate:"max=20"`
	Bankname    string `json:"bankname" validate:"required,max=50"`
	BankAccount string `json:"bankAccount" validate:"required,max=15"`
	Channel     string `json:"channel" validate:"required,max=20" enum:"Google,Youtube,Facebook" example:"Google"`
	TrueWallet  string `json:"trueWallet" validate:"required,max=20"`
	Contact     string `json:"contact" validate:"max=255"`
	Note        string `json:"note" validate:"max=255"`
	Course      string `json:"course" validate:"max=50"`
}

type UserBody struct {
	Fullname string `json:"fullname" validate:"required,8,max=30"`
	// Phone         string   `json:"phone" validate:"required,number,min=10,max=12"`
	Email         string   `json:"email" validate:"required,email"`
	GroupId       *int64   `json:"groupId"`
	Status        string   `json:"status" validate:"required"`
	PermissionIds *[]int64 `json:"permissionIds"`
}

type UserList struct {
	Id           int64      `json:"id"`
	MemberCode   string     `json:"memberCode"`
	Promotion    string     `json:"promotion"`
	Fullname     string     `json:"fullname"`
	Bankname     string     `json:"bankname"`
	BankAccount  string     `json:"bankAccount"`
	Channel      string     `json:"channel"`
	Credit       int        `json:"credit"`
	Ip           string     `json:"ip"`
	IpRegistered string     `json:"ipRegistered"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	LogedinAt    *time.Time `json:"logedinAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type UserDetail struct {
	Id          int64  `json:"id"`
	Partner     string `json:"partner"`
	MemberCode  string `json:"memberCode"`
	Phone       string `json:"phone"`
	Promotion   string `json:"promotion"`
	Fullname    string `json:"fullname"`
	Bankname    string `json:"bankname"`
	BankAccount string `json:"bankAccount"`
	Channel     string `json:"channel"`
	TrueWallet  string `json:"trueWallet"`
	Contact     string `json:"contact"`
	Note        string `json:"note"`
	Course      string `json:"course"`
}

type UserUpdatePassword struct {
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type UserByPhone struct {
	Id    int64  `json:"id"`
	Phone string `json:"phone"`
}

type UserLoginLog struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	Ip        string    `json:"ip"`
	CreatedAt time.Time `json:"createdAt"`
}
