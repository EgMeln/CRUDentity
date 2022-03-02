package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MyMockedAuth struct {
	authUser mocks.Authentication
	token    mocks.Tokens
}

func TestAuthenticationService_SignIn(t *testing.T) {
	testAuth := new(MyMockedAuth)
	user := &model.User{Username: "test", Password: "1234", Admin: false}
	cfg, err := config.New()
	testAuth.token.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(
		func(e context.Context, username string) error {
			return nil
		})
	testAuth.token.On("Add", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.Token")).Return(
		func(e context.Context, token *model.Token) error {
			return nil
		})
	access := NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	serviceAuth := NewAuthService(&testAuth.token, access, refresh, cfg.HashSalt)
	_, _, err = serviceAuth.SignIn(context.Background(), user)
	require.NoError(t, err)
}

func TestAuthenticationService_RefreshToken(t *testing.T) {
	testAuth := new(MyMockedAuth)
	user := &model.User{Username: "test", Password: "1234", Admin: false}
	cfg, err := config.New()
	testAuth.token.On("Delete", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).Return(
		func(e context.Context, username string) error {
			return nil
		})
	testAuth.token.On("Add", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("*model.Token")).Return(
		func(e context.Context, token *model.Token) error {
			return nil
		})
	access := NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	serviceAuth := NewAuthService(&testAuth.token, access, refresh, cfg.HashSalt)
	_, _, err = serviceAuth.SignIn(context.Background(), user)
	require.NoError(t, err)
}
