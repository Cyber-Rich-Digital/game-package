package model

type Success struct {
	Message string `json:"message"`
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
