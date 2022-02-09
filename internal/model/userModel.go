// Package model contain model of struct
package model

import validation "github.com/go-ozzo/ozzo-validation"

const minUsername = 5
const maxUsername = 100

// User struct that contain record info about user
type User struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Admin    bool   `json:"admin" form:"admin"`
}

// Validate checks the validity of the username and password
func (u *User) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Username, validation.Required, validation.Length(minUsername, maxUsername)),
		validation.Field(&u.Password, validation.Required),
		validation.Field(&u.Admin, validation.In("true,false")))
}
