package model

// Token struct that contain record info about token
type Token struct {
	Username     string `json:"username" form:"username"`
	RefreshToken string `json:"token" form:"token"`
}
