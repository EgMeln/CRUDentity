package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type AuthenticationService struct {
	conn         repository.Users
	accessToken  *JWTService
	refreshToken *JWTService
	hashSalt     string
}

func NewAuthenticationServicePostgres(rep *repository.Postgres, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{conn: rep, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}
func NewAuthenticationServiceMongo(rep *repository.Mongo, access, refresh *JWTService, hashSalt string) *AuthenticationService {
	return &AuthenticationService{conn: rep, accessToken: access, refreshToken: refresh, hashSalt: hashSalt}
}

func (srv *AuthenticationService) SignUp(e context.Context, user *model.User) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(srv.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	_, err := srv.conn.GetUser(e, user.Username)
	if err != nil {
		return srv.conn.AddUser(e, user)
	}

	return fmt.Errorf("can't insert user %w", err)
}

func (srv *AuthenticationService) SignIn(e context.Context, user *model.User) (string, string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(srv.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	userSignIn, err := srv.conn.GetUser(e, user.Username)

	if err != nil {
		return "", "", fmt.Errorf("error witn get user %w", err)
	}
	if user.Password != userSignIn.Password {
		return "", "", fmt.Errorf("error witn password %w", err)
	}
	if user.Admin != userSignIn.Admin {
		return "", "", fmt.Errorf("error witn role %w", err)
	}
	accessToken, refreshToken, err := GenerateRefreshAccessToken(srv.accessToken, srv.refreshToken, user)

	shRefreshToken := refreshToken
	pwd = sha1.New()
	pwd.Write([]byte(shRefreshToken))
	pwd.Write([]byte(srv.hashSalt))
	shRefreshToken = fmt.Sprintf("%x", pwd.Sum(nil))

	err = srv.conn.AddToken(e, user.Username, shRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("can't add refresh token %w", err)
	}
	return accessToken, refreshToken, err
}
