package model

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	Id        int            `json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeleteAt  gorm.DeletedAt `json:"deleteAt"`
}

type CreateGroup struct {
	Name string `json:"group" validate:"required"`
}
