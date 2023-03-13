package model

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUser struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=30"`
}

type GetParam struct {
	Id int `uri:"id" binding:"required"`
}

type UserReponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
