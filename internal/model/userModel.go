// Package model contain model of struct
package model

// User struct that contain record info about user
type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}
