package model

import "time"

type Notifications struct {
	Id        int       `json:"id"`
	Total     int       `json:"total" gorm:"default:1"`
	DeviceId  int       `json:"deviceId"`
	TagId     int       `json:"tagId"`
	CreatedAt time.Time `json:"createdAt"`
}
