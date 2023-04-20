package model

type Menu struct {
	Id    int64      `json:"id"`
	Title string     `json:"title"`
	Name  string     `json:"name"`
	Read  bool       `json:"read"`
	Write bool       `json:"write"`
	List  *[]SubMenu `json:"list"`
}

type SubMenu struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
	Read  bool   `json:"read"`
	Write bool   `json:"write"`
}
