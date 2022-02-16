package config

import "github.com/golang-jwt/jwt"

// Claim for JWT
type Claim struct {
	jwt.StandardClaims
	Username string `json:"username" form:"username"`
	Admin    bool   `json:"admin" form:"admin"`
}
