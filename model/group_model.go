package model

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	Id        int64          `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

type CreateGroup struct {
	Name string `json:"name" validate:"required"`
}

type GroupList struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	AdminCount int64  `json:"adminCount"`
}
