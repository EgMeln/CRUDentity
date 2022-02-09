package model

// Token struct that contain record info about token
type Token struct {
	Username     string `json:"username" form:"username"`
	Admin        bool   `json:"admin" form:"admin"`
	RefreshToken string `json:"token" form:"token"`
}
