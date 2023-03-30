package model

type Success struct {
<<<<<<< HEAD
	Message string `json:"message" example:"success" `
}

type SuccessWithToken struct {
	Message string      `json:"message" example:"success" `
	Token   interface{} `json:"token"`
}

type Pagination struct {
	List  interface{} `json:"list" `
	Total int64       `json:"total"`
=======
	Message string `json:"message"`
>>>>>>> c28783323224d6392aa28fde7865115da67ee4ef
}

type SuccessWithData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type SuccessWithList struct {
	Message string      `json:"message"`
	List    interface{} `json:"list"`
}

type SuccessWithPagination struct {
	Message string      `json:"message"`
	List    interface{} `json:"list"`
	Total   int64       `json:"total"`
}
