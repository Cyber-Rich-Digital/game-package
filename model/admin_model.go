package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	Id           int64          `json:"id"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Fullname     string         `json:"fullname"`
	Firstname    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Role         string         `json:"role"`
	Status       string         `json:"status"`
	AdminGroupId int64          `json:"adminGroupId"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt"`
	LogedinAt    *time.Time     `json:"logedinAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type CreateAdmin struct {
	Username      string   `json:"username" validate:"required,8,max=30"`
	Password      string   `json:"password" validate:"required,8,max=30"`
	Fullname      string   `json:"fullname" validate:"required,8,max=30"`
	Phone         string   `json:"phone" validate:"required,min=10,max=12"`
	Email         string   `json:"email" validate:"required,email"`
	RoleId        string   `json:"roleId" validate:"required"`
	Status        string   `json:"status" validate:"required"`
	AdminGroupId  int64    `json:"adminGroupId" validate:"required"`
	PermissionIds *[]int64 `json:"permissionIds" validate:"required"`
}

type LoginAdmin struct {
	Username string `json:"username" validate:"required,min=8,max=30"`
	Password string `json:"password" validate:"required,8,max=30"`
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

type AdminUpdateGroup struct {
	GroupId       int64   `json:"groupId" validate:"required"`
	PermissionIds []int64 `json:"permissionId" validate:"required"`
}

type AdminGroupQuery struct {
	Page  int `form:"page" validate:"required,min=1"`
	Limit int `form:"limit" validate:"required,min=1,max=100"`
}

type AdminListQuery struct {
	Page   int    `form:"page" validate:"required,min=1"`
	Limit  int    `form:"limit" validate:"required,min=1,max=100"`
	Search string `form:"search"`
	Status string `form:"status"`
}

type AdminGroupPagination struct {
	Total int64            `json:"total"`
	List  []GroupCountList `json:"list"`
}

type UpdateAdmin struct {
	Fullname     string `json:"fullname" validate:"required,8,max=30"`
	Firstname    string `json:"firstname" validate:"required,8,max=30"`
	Lastname     string `json:"lastname" validate:"required,8,max=30"`
	Phone        string `json:"phone" validate:"required,min=10,max=12"`
	Email        string `json:"email" validate:"required,email"`
	Role         string `json:"roleId" validate:"required"`
	Status       string `json:"status" validate:"required"`
	AdminGroupId *int64 `json:"adminGroupId" validate:"required"`
}

type AdminBody struct {
	Fullname string `json:"fullname" validate:"required,8,max=30"`
	// Phone         string   `json:"phone" validate:"required,number,min=10,max=12"`
	Email         string   `json:"email" validate:"required,email"`
	GroupId       *int64   `json:"groupId"`
	Status        string   `json:"status" validate:"required"`
	PermissionIds *[]int64 `json:"permissionIds"`
}

type AdminPermission struct {
	AdminId      int64 `json:"adminId"`
	PermissionId int64 `json:"permissionId"`
}

type AdminList struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type AdminDetail struct {
	Id             int64            `json:"id"`
	Username       string           `json:"username"`
	Fullname       string           `json:"fullname"`
	Phone          string           `json:"phone"`
	Email          string           `json:"email"`
	Role           string           `json:"role"`
	Status         string           `json:"status"`
	PermissionList []PermissionList `json:"permissionList"`
	Group          *GroupDetail     `json:"group"`
}

type AdminGroupId struct {
	AdminGroupId int `json:"adminGroupId"`
}

type AdminUpdatePassword struct {
	Password string `json:"password" validate:"required,min=8,max=30"`
}
