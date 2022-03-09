// Package service business logic
package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/golang-jwt/jwt"
)

// AuthenticationService struct for rep
type AuthenticationService struct {
	token        repository.Tokens
	accessToken  *JWTService
	refreshToken *JWTService
	hashSalt     string
}

// JWTService struct for fields of a token
type JWTService struct {
	Key           []byte
	TokenLifeTime time.Duration
}

// NewAuthServicePostgres used for setting postgres services
func NewAuthServicePostgres(rep2 *repository.PostgresToken, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{token: rep2, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}

// NewAuthServiceMongo used for setting mongo services
func NewAuthServiceMongo(rep2 *repository.MongoToken, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{token: rep2, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}

// NewAuthService used for setting mongo services
func NewAuthService(rep2 repository.Tokens, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{token: rep2, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}

// NewJWTService used for token setting mongo services
func NewJWTService(key []byte, tokenLifeTime time.Duration) *JWTService {
	return &JWTService{Key: key, TokenLifeTime: tokenLifeTime}
}

// SignIn generate token
func (srv *AuthenticationService) SignIn(e context.Context, user *model.User) (access, refresh string, ok error) {
	_ = srv.token.Delete(e, user.Username)

	accessToken, refreshToken, err := srv.GenerateRefreshAccessToken(srv.accessToken, srv.refreshToken, user)
	if err != nil {
		return "", "", fmt.Errorf("can't generate tokens %w", err)
	}

	shRefreshToken := refreshToken
	pwd := sha256.New()
	pwd.Write([]byte(shRefreshToken))
	pwd.Write([]byte(srv.hashSalt))
	shRefreshToken = fmt.Sprintf("%x", pwd.Sum(nil))

	err = srv.token.Add(e, &model.Token{Username: user.Username, RefreshToken: shRefreshToken})
	if err != nil {
		return "", "", fmt.Errorf("can't add refresh token %w", err)
	}
	return accessToken, refreshToken, err
}

// GenerateToken generate token
func (srv *AuthenticationService) GenerateToken(jwtService *JWTService, user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &config.Claim{
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

// GenerateRefreshAccessToken generating two tokens
func (srv *AuthenticationService) GenerateRefreshAccessToken(access, refresh *JWTService, user *model.User) (accessToken, refreshToken string, err error) {
	accessToken, err = srv.GenerateToken(access, user)
	if err != nil {
		return "", "", fmt.Errorf("error with acces token %w", err)
	}
	refreshToken, err = srv.GenerateToken(refresh, user)
	if err != nil {
		return "", "", fmt.Errorf("error with refresh token %w", err)
	}

	if err != nil {
		return "", "", fmt.Errorf("can't add token %w", err)
	}
	return accessToken, refreshToken, nil
}

// RefreshToken generating a new refresh token
func (srv *AuthenticationService) RefreshToken(e context.Context, user *model.User) (accessToken, refreshToken string, err error) {
	err = srv.token.Delete(e, user.Username)
	if err != nil {
		return "", "", fmt.Errorf("can't delete token %w", err)
	}
	accessToken, refreshToken, err = srv.GenerateRefreshAccessToken(srv.accessToken, srv.refreshToken, user)
	if err != nil {
		return "", "", fmt.Errorf("can't generate tokens %w", err)
	}
	shRefreshToken := refreshToken
	pwd := sha256.New()
	pwd.Write([]byte(shRefreshToken))
	pwd.Write([]byte(srv.hashSalt))
	shRefreshToken = fmt.Sprintf("%x", pwd.Sum(nil))

	err = srv.token.Add(e, &model.Token{Username: user.Username, RefreshToken: shRefreshToken})
	if err != nil {
		return "", "", fmt.Errorf("can't add refresh token %w", err)
	}
	return accessToken, refreshToken, nil
}
