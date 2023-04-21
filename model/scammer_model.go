package model

import (
	"time"
)

type Scammer struct {
	Id          int64      `json:"id"`
	Fullname    *string    `json:"fullname"`
	Firstname   *string    `json:"firstname"`
	Lastname    *string    `json:"lastname"`
	Bankname    *string    `json:"bankname"`
	BankAccount *string    `json:"bankAccount"`
	Phone       *string    `json:"phone"`
	Reason      *string    `json:"reason"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type CreateScammer struct {
	Fullname    *string `json:"fullname"`
	Bankname    *string `json:"bankname" validate:"max=50"`
	BankAccount *string `json:"bankAccount" validate:"max=15"`
	Phone       *string `json:"phone" validate:"required,min=10,max=12"`
	Reason      *string `json:"reason" validate:"max=255"`
}

type ScammertList struct {
	Id          int64      `json:"id"`
	Fullname    *string    `json:"fullname"`
	Bankname    *string    `json:"bankname"`
	BankAccount *string    `json:"bankAccount"`
	Phone       *string    `json:"phone"`
	Reason      *string    `json:"reason"`
	CreatedAt   *time.Time `json:"createdAt"`
}

type ScammerDetail struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type ScammerQuery struct {
	Page      int     `form:"page" default:"1"`
	Limit     int     `form:"limit" default:"10"`
	DateStart string  `form:"dateStart" default:"2020-01-01 00:00:00"`
	DateEnd   string  `form:"dateEnd" default:"2020-01-01 00:00:00"`
	BankName  *string `form:"bankName" default:""`
	Filter    *string `form:"filter" default:""`
}
