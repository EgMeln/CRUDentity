package model

type Token struct {
	Username     string `json:"username" form:"username"`
	Admin        bool   `json:"admin" form:"admin"`
	RefreshToken string `json:"token" form:"token"`
}
