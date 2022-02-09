// Package handlers contain function for handling request
package handlers

import (
	"net/http"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/request"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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
	user, err = handler.service.Get(e.Request().Context(), username)
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
	err := handler.service.Update(e.Request().Context(), &model.User{Username: username, Password: c.Password, Admin: c.Admin})
	if err != nil {
		log.Warnf("Update user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting user
func (handler *UserHandler) Delete(e echo.Context) error {
	username := e.Param("username")
	err := handler.service.Delete(e.Request().Context(), username)
	if err != nil {
		log.Warnf("Delete user error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, e)
	}
	return e.JSON(http.StatusOK, e)
}
