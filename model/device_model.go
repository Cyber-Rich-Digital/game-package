package model

import "time"

type Device struct {
	Id         int       `json:"id"`
	Os         string    `json:"os"`
	Version    string    `json:"version"`
	FcmToken   string    `json:"fcmToken"`
	HardwareId string    `json:"hardwareId"`
	WebsiteId  int       `json:"websiteId"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:CURRENT_TIMESTAMP"`
}

type DeviceBody struct {
	FcmToken   string `json:"fcmToken"`
	HardwareId string `json:"hardwareId"`
	WebsiteId  int    `json:"websiteId" validate:"required"`
}
