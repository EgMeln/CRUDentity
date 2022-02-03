package service

import (
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/golang-jwt/jwt"
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

	if err != nil {
		return "", "", fmt.Errorf("can't add token %w", err)
	}
	return accessToken, refreshToken, nil
}
