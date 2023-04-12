package model

import (
	"time"
)

type Recomment struct {
	Id        int64      `json:"id"`
	Title     *string    `json:"title"`
	Status    *string    `json:"status"`
	Url       *string    `json:"url"`
	CreatedAt *time.Time `json:"createdAt"`
}

type CreateRecomment struct {
	Title  *string `json:"title" validate:"required,max=255"`
	Status *string `json:"status" validate:"required,max=10" enums:"ACTIVE,DEACTIVE" select:"ACTIVE,DEACTIVE"`
	Url    *string `json:"url" validate:"required,max=255" example:"https://www.facebook.com/.../"`
}

type RecommentList struct {
	Id        int64      `json:"id"`
	Title     *string    `json:"title"`
	Status    *string    `json:"status"`
	Url       *string    `json:"url"`
	CreatedAt *time.Time `json:"createdAt"`
}

type RecommentQuery struct {
	Page   int    `form:"page" validate:"required,min=1" example:"1"`
	Limit  int    `form:"limit" validate:"required,min=1,max=100" example:"10"`
	Status string `form:"status" example:""`
	Filter string `form:"filter" example:""`
}
