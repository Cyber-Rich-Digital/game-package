package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        uint           `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

type CreateUser struct {
	Username string `json:"username" validate:"required,min=6,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type CreateAdmin struct {
	Username string `json:"username" validate:"required,min=6,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type Login struct {
	Email    string `json:"email" validate:"required,min=6,max=30"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type UserJwt struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UserQuery struct {
	Page   int    `form:"page" default:"1" min:"1"`
	Limit  int    `form:"limit" default:"10" min:"1" max:"100"`
	Search string `form:"search"`
	Sort   int    `form:"sort"`
}

type UserResponse struct {
	Id        uint          `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"createdAt"`
	WebTotal  int           `json:"webTotal"`
	Websites  []WebsiteList `json:"websites" gorm:"foreignKey:UserId;references:Id"`
}

type UserListResponse struct {
	Id        uint          `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email"`
	CreatedAt time.Time     `json:"createdAt"`
	WebTotal  int           `json:"webTotal"`
	Websites  []WebsiteList `json:"websites"`
}

type UserAdminResponse struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserChangePassword struct {
	Password    string `json:"password" validate:"required,min=6,max=30"`
	OldPassword string `json:"oldPassword" validate:"required,min=6,max=30"`
}

type AdminChangePassword struct {
	Password string `json:"password" validate:"required,min=6,max=30"`
}
