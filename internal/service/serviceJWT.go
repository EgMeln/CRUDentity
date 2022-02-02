package service

import (
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type JWTService struct {
	Key           []byte
	TokenLifeTime time.Duration
}

func NewJWTService(key []byte, tokenLifeTime time.Duration) *JWTService {
	return &JWTService{Key: key, TokenLifeTime: tokenLifeTime}
}

func GenerateToken(jwtService *JWTService, user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtService.TokenLifeTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		Admin:    user.Admin,
	})
	tokenString, err := token.SignedString(jwtService.Key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshAccessToken(access *JWTService, refresh *JWTService, user *model.User) (string, string, error) {
	accessToken, err := GenerateToken(access, user)
	if err != nil {
		return "", "", fmt.Errorf("error witn acces token %w", err)
	}
	refreshToken, err := GenerateToken(refresh, user)

	if err != nil {
		return "", "", fmt.Errorf("error witn refresh token %w", err)
	}

	return accessToken, refreshToken, nil
}

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

func TokenRefresh(access, refresh *JWTService) echo.MiddlewareFunc {
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
						_, _, _ = GenerateRefreshAccessToken(access, refresh, &model.User{
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
