package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Id        int       `json:"id"`
	WebsiteId int       `json:"websiteId"`
	Name      string    `json:"name"`
	CreateAt  time.Time `json:"createdAt"`
	DeletedAt gorm.DeletedAt
}

type TagParam struct {
	WebsiteId int `uri:"website_id" binding:"required"`
	DeviceId  float64
}

type TagInsert struct {
	Id        int
	WebsiteId int    `json:"websiteId" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

type TagResponse struct {
	Id        int       `json:"id" gorm:"primary_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Total     int       `json:"total"`
}

type TagList struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	WebsiteId int    `json:"websiteId"`
}
