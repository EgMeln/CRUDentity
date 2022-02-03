package middlewares

import (
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func TokenRefresh(access, refresh *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get("user") == nil {
				return next(c)
			}
			u := c.Get("user").(*jwt.Token)

			claims := u.Claims.(*model.Claim)

			if time.Until(time.Unix(claims.ExpiresAt, 0)) < 6*time.Minute {
				if refresh.Key != nil {
					// Parses token and checks if it valid.
					tkn, err := jwt.ParseWithClaims(string(refresh.Key), claims, func(token *jwt.Token) (interface{}, error) {
						return refresh.Key, nil
					})
					if err != nil {
						if err == jwt.ErrSignatureInvalid {
							c.Response().Writer.WriteHeader(http.StatusUnauthorized)
						}
					}

					if tkn != nil && tkn.Valid {
						// If everything is good, update tokens.
						_, _, _ = service.GenerateRefreshAccessToken(access, refresh, &model.User{
							Username: claims.Username,
							Admin:    claims.Admin,
						})
					}
				}
			}
			return next(c)
		}
	}
}
