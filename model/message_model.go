package model

import "time"

type Message struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	WebsiteID int       `json:"websiteId"`
	TagID     int       `json:"tagId"`
	CreateAt  time.Time `json:"createdAt"`
}

type MessageParam struct {
	TagId int `uri:"tag_id" binding:"required"`
}

type MessageBody struct {
	Message string `json:"message" validate:"required"`
	ApiKey  string `json:"apiKey" validate:"required"`
	Tag     string `json:"tag" validate:"required"`
}

type MessageCreate struct {
	Message   string `json:"message" validate:"required"`
	WebsiteId int    `json:"websiteId" validate:"required"`
	TagId     int    `json:"tagId" validate:"required"`
}

type MessageResponse struct {
	Id        int       `json:"id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

type MessageRead struct {
	TagId    int `json:"tagId" validate:"required"`
	DeviceId int `json:"deviceId" validate:"required"`
}
