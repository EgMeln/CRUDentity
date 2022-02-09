package handlers

import (
	"fmt"
	"net/http"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/request"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// AuthenticationHandler struct that contain repository linc
type AuthenticationHandler struct {
	service *service.AuthenticationService
}

// NewServiceAuthentication add new authentication handler
func NewServiceAuthentication(srv *service.AuthenticationService) AuthenticationHandler {
	return AuthenticationHandler{service: srv}
}

// SignUp user
func (handler *AuthenticationHandler) SignUp(e echo.Context) error { //nolint:dupl //Different business logic
	user := new(request.SignUpUser)
	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err := handler.service.SignUp(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password})
	if err != nil {
		log.Warnf("Sign up error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}
	return e.JSON(http.StatusOK, user)
}

// SignIn user and generate token
func (handler *AuthenticationHandler) SignIn(e echo.Context) error {
	user := new(request.SignInUser)

	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	accessToken, refreshToken, err := handler.service.SignIn(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password, Admin: user.Admin})

	if err != nil {
		log.Warnf("Sign in error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}

	return e.JSON(http.StatusOK, fmt.Sprintf("access token: %s, refresh token: %s", accessToken, refreshToken))
}

// Refresh token
func (handler *AuthenticationHandler) Refresh(e echo.Context) error {
	user := new(request.RefreshToken)
	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	accessToken, refreshToken, err := handler.service.RefreshToken(e.Request().Context(), &model.User{Username: user.Username, Admin: user.Admin})
	if err != nil {
		log.Warnf("Refresh token error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}

	return e.JSON(http.StatusOK, fmt.Sprintf("new access token: %s, new refresh token: %s", accessToken, refreshToken))
}
