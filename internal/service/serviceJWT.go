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
	token := jwt.NewWithClaims(jwt.SigningMethodES256, &model.Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtService.TokenLifeTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Username: user.Username,
		Admin:    user.Admin,
	})
	tokenString, err := token.SignedString(jwtService.Key)
	if err != nil {
		return "", fmt.Errorf("can't generate token %w", err)
	}
	return tokenString, nil
}
func GenerateRefreshAccessToken(access *JWTService, refresh *JWTService, user *model.User) (string, string, error) {
	accessToke, err := GenerateToken(access, user)
	if err != nil {
		return "", "", fmt.Errorf("can't create access token %w", err)
	}
	refreshToken, err := GenerateToken(refresh, user)
	if err != nil {
		return "", "", fmt.Errorf("can't create refresh token %w", err)
	}
	return accessToke, refreshToken, nil
}

func (JWTService *JWTService) ParseToken(accessToken string) (*model.Claim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &model.Claim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return JWTService.Key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := token.Claims.(*model.Claim)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
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
