package middlewares

import (
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		if e.Get("user") == nil {
			return next(e)
		}
		u := e.Get("user").(*jwt.Token)
		claims := u.Claims.(*model.Claim)

		if !claims.Admin {
			return echo.NewHTTPError(http.StatusNotAcceptable, "have no access")
		}
		return next(e)
	}
}
