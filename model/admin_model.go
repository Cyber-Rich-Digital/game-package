package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	Id        int64          `json:"id"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	Fullname  string         `json:"fullname"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	LogedinAt *time.Time     `json:"logedinAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type CreateAdmin struct {
	Username    string       `json:"username" validate:"required,min=6,max=30"`
	Password    string       `json:"password" validate:"required,min=6,max=30"`
	Fullname    string       `json:"fullname" validate:"required,min=6,max=30"`
	Phone       string       `json:"phone" validate:"required,min=10,max=12"`
	Email       string       `json:"email" validate:"required,email"`
	RoleId      string       `json:"roleId" validate:"required"`
	Status      string       `json:"status" validate:"required"`
	Permissions []Permission `json:"permissions" validate:"required"`
}

type LoginAdmin struct {
	Username string `json:"username" validate:"required,min=8,max=30"`
	Password string `json:"password" validate:"required,min=6,max=30"`
	IP       string `json:"ip"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type AdminLoginUpdate struct {
	IP        string    `json:"ip"`
	LogedinAt time.Time `json:"logedinAt"`
}

type AdminCreateGroup struct {
	GroupId       int64   `json:"groupId" validate:"required"`
	PermissionIds []int64 `json:"permissionIds" validate:"required"`
}

type AdminPermissionList struct {
	GroupId      int64 `json:"groupId"`
	PermissionId int64 `json:"permissionId"`
}

type AdminGroupPermission struct {
	Id           int64     `json:"id"`
	GroupId      int64     `json:"groupId"`
	PermissionId int64     `json:"permissionId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	DeletedAt    time.Time `json:"deletedAt"`
}

type AdminGroupPermissionResponse struct {
	Id          int64            `json:"id"`
	Name        string           `json:"name"`
	Permissions []PermissionList `json:"permissions"`
}