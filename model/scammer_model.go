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
	Page      int        `form:"page" example:"1"`
	Limit     int        `form:"limit" example:"10"`
	DateStart *time.Time `form:"dateStart" example:"2020-01-01 00:00:00"`
	DateEnd   *time.Time `form:"dateEnd" example:"2020-01-01 00:00:00"`
	BankName  *string    `form:"bankName" example:"-"`
	Filter    *string    `form:"filter" example:""`
}
