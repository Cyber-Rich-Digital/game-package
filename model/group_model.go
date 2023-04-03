package model

import (
	"time"
)

type Group struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type CreateGroup struct {
	Name string `json:"name" validate:"required"`
}

type GroupCountList struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AdminCount int64  `json:"adminCount"`
}

type GroupDetail struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DeleteGroup struct {
	Id int64 `json:"id" validate:"required"`
}
