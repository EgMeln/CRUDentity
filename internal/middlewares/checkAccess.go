// Package middlewares for JWT
package middlewares

import (
	"net/http"

	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// CheckAccess for checking roll Admin
func CheckAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if e.Get("user") == nil {
			return next(e)
		}
		u := e.Get("user").(*jwt.Token)
		claims := u.Claims.(*config.Claim)

		if !claims.Admin {
			return echo.NewHTTPError(http.StatusNotAcceptable, "have no access")
		}
		return next(e)
	}
}
