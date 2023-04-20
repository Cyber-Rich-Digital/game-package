package model

import (
	"time"
)

type Permission struct {
	Id            int64      `json:"id"`
	PermissionKey string     `json:"permissionKey"`
	Name          string     `json:"name"`
	Main          bool       `json:"main"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deleteAt"`
}

type CreatePermission struct {
	Permissions []PermissionName `json:"permissions" validate:"required"`
}

type PermissionName struct {
	Name          string `json:"name" validate:"required"`
	PermissionKey string `json:"permissionKey" validate:"required"`
	Main          bool   `json:"isMain" default:"false"`
}

type PermissionList struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	IsRead  bool   `json:"read"`
	IsWrite bool   `json:"write"`
}

type PermissionObj struct {
	Id    int64 `json:"id"`
	Read  bool  `json:"read" default:"false"`
	Write bool  `json:"write" default:"false"`
}

type DeletePermission struct {
	PermissionIds []int64 `json:"permissionIds" validate:"required"`
}
