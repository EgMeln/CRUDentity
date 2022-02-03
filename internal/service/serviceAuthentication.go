package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthenticationService struct {
	conn         repository.Users
	token        repository.Tokens
	accessToken  *JWTService
	refreshToken *JWTService
	hashSalt     string
}

type JWTService struct {
	Key           []byte
	TokenLifeTime time.Duration
}

func NewAuthenticationServicePostgres(rep *repository.PostgresUser, rep2 *repository.PostgresToken, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{conn: rep, token: rep2, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}
func NewAuthenticationServiceMongo(rep *repository.MongoUser, rep2 *repository.MongoToken, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{conn: rep, token: rep2, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}
func NewJWTService(key []byte, tokenLifeTime time.Duration) *JWTService {
	return &JWTService{Key: key, TokenLifeTime: tokenLifeTime}
}

func (srv *AuthenticationService) SignUp(e context.Context, user *model.User) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(srv.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	_, err := srv.conn.Get(e, user.Username)
	if err != nil {
		return srv.conn.Add(e, user)
	}

	return fmt.Errorf("can't insert user %w", err)
}

func (srv *AuthenticationService) SignIn(e context.Context, user *model.User) (string, string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(srv.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	userSignIn, err := srv.conn.Get(e, user.Username)

	if err != nil {
		return "", "", fmt.Errorf("error witn get user %w", err)
	}
	if user.Password != userSignIn.Password {
		return "", "", fmt.Errorf("error witn password %w", err)
	}
	if user.Admin != userSignIn.Admin {
		return "", "", fmt.Errorf("error witn role %w", err)
	}
	accessToken, refreshToken, err := srv.GenerateRefreshAccessToken(srv.accessToken, srv.refreshToken, user)

	shRefreshToken := refreshToken
	pwd = sha1.New()
	pwd.Write([]byte(shRefreshToken))
	pwd.Write([]byte(srv.hashSalt))
	shRefreshToken = fmt.Sprintf("%x", pwd.Sum(nil))

	err = srv.token.Add(e, &model.Token{Username: userSignIn.Username, Admin: userSignIn.Admin, RefreshToken: shRefreshToken})
	if err != nil {
		return "", "", fmt.Errorf("can't add refresh token %w", err)
	}
	return accessToken, refreshToken, err
}

func (srv *AuthenticationService) GenerateToken(jwtService *JWTService, user *model.User) (string, error) {
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
func (srv *AuthenticationService) GenerateRefreshAccessToken(access *JWTService, refresh *JWTService, user *model.User) (string, string, error) {
	accessToken, err := srv.GenerateToken(access, user)
	if err != nil {
		return "", "", fmt.Errorf("error witn acces token %w", err)
	}
	refreshToken, err := srv.GenerateToken(refresh, user)
	if err != nil {
		return "", "", fmt.Errorf("error witn refresh token %w", err)
	}

	if err != nil {
		return "", "", fmt.Errorf("can't add token %w", err)
	}
	return accessToken, refreshToken, nil
}
func (srv *AuthenticationService) RefreshToken(e context.Context, user *model.User) (string, string, error) {

	err := srv.token.Delete(e, user.Username)
	if err != nil {
		return "", "", fmt.Errorf("can't delete token %w", err)
	}
	accessToken, refreshToken, err := srv.GenerateRefreshAccessToken(srv.accessToken, srv.refreshToken, user)

	shRefreshToken := refreshToken
	pwd := sha1.New()
	pwd.Write([]byte(shRefreshToken))
	pwd.Write([]byte(srv.hashSalt))
	shRefreshToken = fmt.Sprintf("%x", pwd.Sum(nil))

	err = srv.token.Add(e, &model.Token{Username: user.Username, Admin: user.Admin, RefreshToken: shRefreshToken})
	if err != nil {
		return "", "", fmt.Errorf("can't add refresh token %w", err)
	}
	return accessToken, refreshToken, nil
}
