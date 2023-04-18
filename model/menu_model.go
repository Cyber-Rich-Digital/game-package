package model

type Menu struct {
	Id    int64      `json:"id"`
	Title string     `json:"title"`
	Name  string     `json:"name"`
	View  bool       `json:"view"`
	Edit  bool       `json:"edit"`
	List  *[]SubMenu `json:"list"`
}

type SubMenu struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Name  string `json:"name"`
	View  bool   `json:"view"`
	Edit  bool   `json:"edit"`
}
