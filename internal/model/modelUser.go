package model

type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
	Token    string `json:"token" form:"token"`
}
