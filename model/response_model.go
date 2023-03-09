package model

type Success struct {
	Message string `json:"message" example:"success" `
}

type SuccessWithToken struct {
	Message string      `json:"message" example:"success" `
	Token   interface{} `json:"token"`
}

type Pagination struct {
	List  interface{} `json:"list"`
	Total int         `json:"total"`
}

type Error struct {
	Message string `json:"message" example:"error" `
}
