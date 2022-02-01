package model

import "github.com/golang-jwt/jwt"

type Claim struct {
	jwt.StandardClaims
	Username string `json:"username" form:"username"`
	Admin    bool   `json:"admin" form:"admin"`
}
