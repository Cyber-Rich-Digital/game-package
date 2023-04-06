package model

import (
	"time"
)

type Permission struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeleteAt  time.Time `json:"deleteAt"`
}

type CreatePermission struct {
	Permissions []PermissionName `json:"permissions" validate:"required"`
}

type PermissionName struct {
	Name string `json:"name" validate:"required"`
}

type PermissionList struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DeletePermission struct {
	PermissionIds []int64 `json:"permissionIds" validate:"required"`
}
