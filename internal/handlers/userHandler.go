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
// @Summary sign-in
// @ID sign-in-user
// @Produce json
// @Param request body request.SignInUser true "sign in user"
// @Success 200 {object} request.SignInUser
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /auth/sign-in [post]
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

	userSignIn, err := handler.userService.Get(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password})

	if err != nil {
		log.Warnf("User doesn't exist: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}

	accessToken, refreshToken, err := handler.authService.SignIn(e.Request().Context(), userSignIn)

	if err != nil {
		log.Warnf("Tokens generate error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, user)
	}
	return e.JSON(http.StatusOK, fmt.Sprintf("access token: %s, refresh token: %s", accessToken, refreshToken))
}

// Add record about user
// @Summary add user
// @ID add-user
// @Produce json
// @Param request body request.SignUpUser true "sign up user"
// @Success 200 {object} request.SignUpUser
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /auth/sign-up [post]
func (handler *UserHandler) Add(e echo.Context) error { //nolint:dupl //Different business logic
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
// @Summary gets all users
// @ID get-all-users
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} model.User
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/users [get]
func (handler *UserHandler) GetAll(e echo.Context) error {
	users, err := handler.userService.GetAll(e.Request().Context())
	if err != nil {
		log.Warnf("Get all users error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, users)
	}
	return e.JSON(http.StatusOK, users)
}

// Get getting user by username
// @Summary get user by username
// @ID get-user-by-username
// @Security ApiKeyAuth
// @Produce json
// @Param username path string true "get user"
// @Success 200 {object} request.GetUser
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/users/{username} [get]
func (handler *UserHandler) Get(e echo.Context) error {
	username := e.Param("username")
	var err error
	getUser, err := handler.userService.Get(e.Request().Context(), &model.User{Username: username})
	if err != nil {
		log.Warnf("Get user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, getUser)
	}
	return e.JSON(http.StatusOK, getUser)
}

// Update updating user
// @Summary update user by username
// @ID update-user-by-username
// @Security ApiKeyAuth
// @Produce json
// @Param request body request.UpdateUser true "update user"
// @Success 200 {object} request.UpdateUser
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/users [put]
func (handler *UserHandler) Update(e echo.Context) error { //nolint:dupl //Different business logic
	c := new(request.UpdateUser)
	if err := e.Bind(c); err != nil {
		log.Warnf("Bind fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := e.Validate(c); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err := handler.userService.Update(e.Request().Context(), &model.User{Username: c.Username, Password: c.Password})
	if err != nil {
		log.Warnf("Update user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting user
// @Summary delete user by username
// @ID delete-user-by-username
// @Security ApiKeyAuth
// @Produce json
// @Param username path string true "update user"
// @Success 200 {string} echo.Context
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/users/{username} [delete]
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
// @Summary refresh token
// @ID refresh-token
// @Security ApiKeyAuth
// @Produce json
// @Param request body request.RefreshToken true "refresh token"
// @Success 200 {string} string
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /user/refresh [post]
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
