// Package handlers contain function for handling request
package handlers

import (
	"net/http"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/request"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
)

// UserHandler struct that contain repository linc
type UserHandler struct {
	service *service.UserService
}

// NewServiceUser add new authentication handler
func NewServiceUser(srv *service.UserService) UserHandler {
	return UserHandler{service: srv}
}

// GetAll getting all users
func (handler *UserHandler) GetAll(e echo.Context) error {
	users, err := handler.service.GetAll(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusBadRequest, users)
	}
	return e.JSON(http.StatusOK, users)
}

// Get getting parking lot by username
func (handler *UserHandler) Get(e echo.Context) error {
	username := e.Param("username")
	var user *model.User
	var err error
	user, err = handler.service.Get(e.Request().Context(), username)
	if err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}
	return e.JSON(http.StatusOK, user)
}

// Update updating user
func (handler *UserHandler) Update(e echo.Context) error {
	username := e.Param("username")
	c := new(request.UpdateUser)
	if err := e.Bind(c); err != nil {
		return err
	}
	err := handler.service.Update(e.Request().Context(), username, c.Password, c.Admin)
	if err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting user
func (handler *UserHandler) Delete(e echo.Context) error {
	username := e.Param("username")
	var err error
	if err != nil {
		return err
	}
	err = handler.service.Delete(e.Request().Context(), username)
	if err != nil {
		return e.JSON(http.StatusBadRequest, e)
	}
	return e.JSON(http.StatusOK, e)
}
