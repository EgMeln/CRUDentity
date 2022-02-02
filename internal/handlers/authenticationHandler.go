package handlers

import (
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthenticationHandler struct {
	service *service.AuthenticationService
}

func NewServiceAuthentication(srv *service.AuthenticationService) AuthenticationHandler {
	return AuthenticationHandler{service: srv}
}

func (handler *AuthenticationHandler) SignUp(e echo.Context) error {
	user := new(model.User)
	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}
	err := handler.service.SignUp(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password, Admin: user.Admin})
	if err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}
	return e.JSON(http.StatusOK, user)
}

func (handler *AuthenticationHandler) SignIn(e echo.Context) error {
	user := new(model.User)

	if err := e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}

	accessToken, refreshToken, err := handler.service.SignIn(e.Request().Context(), user)

	if err != nil {
		return e.JSON(http.StatusBadRequest, fmt.Sprintf("error with tokens"))
	}

	return e.JSON(http.StatusOK, fmt.Sprintf("access token %s, refresh token %s", accessToken, refreshToken))
}
