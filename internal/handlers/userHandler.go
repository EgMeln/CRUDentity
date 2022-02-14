// Package handlers contain function for handling request
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

// UserHandler struct that contain repository linc
type UserHandler struct {
	userService *service.UserService
	authService *service.AuthenticationService
}

// NewServiceUser add new authentication handler
func NewServiceUser(srvUser *service.UserService, srvAuth *service.AuthenticationService) UserHandler {
	return UserHandler{userService: srvUser, authService: srvAuth}
}

// SignIn generate token
func (handler *UserHandler) SignIn(e echo.Context) error {
	user := new(request.SignInUser)

	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := e.Validate(user); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userSignIn, err := handler.userService.Get(e.Request().Context(), user.Username)

	if err != nil {
		log.Warnf("User doesn't exist: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}

	accessToken, refreshToken, err := handler.authService.SignIn(e.Request().Context(), userSignIn, user.Password)

	if err != nil {
		log.Warnf("Tokens generate error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}
	return e.JSON(http.StatusOK, fmt.Sprintf("access token: %s, refresh token: %s", accessToken, refreshToken))
}

// Add record about user
func (handler *UserHandler) Add(e echo.Context) error {
	user := new(request.SignUpUser)
	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := e.Validate(user); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err := handler.userService.Add(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password})
	if err != nil {
		log.Warnf("Sign up error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}
	return e.JSON(http.StatusOK, user)
}

// GetAll getting all users
func (handler *UserHandler) GetAll(e echo.Context) error {
	users, err := handler.userService.GetAll(e.Request().Context())
	if err != nil {
		log.Warnf("Get all users error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, users)
	}
	return e.JSON(http.StatusOK, users)
}

// Get getting parking lot by username
func (handler *UserHandler) Get(e echo.Context) error {
	username := e.Param("username")
	var user *model.User
	var err error
	user, err = handler.userService.Get(e.Request().Context(), username)
	if err != nil {
		log.Warnf("Get user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}
	return e.JSON(http.StatusOK, user)
}

// Update updating user
func (handler *UserHandler) Update(e echo.Context) error {
	username := e.Param("username")
	c := new(request.UpdateUser)
	if err := e.Bind(c); err != nil {
		log.Warnf("Bind fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := e.Validate(c); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err := handler.userService.Update(e.Request().Context(), &model.User{Username: username, Password: c.Password})
	if err != nil {
		log.Warnf("Update user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting user
func (handler *UserHandler) Delete(e echo.Context) error {
	username := e.Param("username")
	err := handler.userService.Delete(e.Request().Context(), username)
	if err != nil {
		log.Warnf("Delete user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, e)
	}
	return e.JSON(http.StatusOK, e)
}

// Refresh generating a new refresh token
func (handler *UserHandler) Refresh(e echo.Context) error {
	user := new(request.RefreshToken)
	if err := e.Bind(user); err != nil {
		log.Warnf("Bind fail: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := e.Validate(user); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	accessToken, refreshToken, err := handler.authService.RefreshToken(e.Request().Context(), &model.User{Username: user.Username})
	if err != nil {
		log.Warnf("Refresh token error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}

	return e.JSON(http.StatusOK, fmt.Sprintf("new access token: %s, new refresh token: %s", accessToken, refreshToken))
}
