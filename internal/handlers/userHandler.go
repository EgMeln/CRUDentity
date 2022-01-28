package handlers

import (
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
}

func NewServiceUser(srv *service.UserService) UserHandler {
	return UserHandler{service: srv}
}
func (handler *UserHandler) Add(e echo.Context) (err error) {
	user := new(model.User)
	if err = e.Bind(user); err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}
	err = handler.service.Add(e.Request().Context(), &model.User{Username: user.Username, Password: user.Password, Admin: user.Admin})
	if err != nil {
		return e.JSON(http.StatusBadRequest, user)
	}
	return e.JSON(http.StatusOK, user)
}
func (handler *UserHandler) GetAll(e echo.Context) error {
	users, err := handler.service.GetAll(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusBadRequest, users)
	}
	return e.JSON(http.StatusOK, users)
}
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
func (handler *UserHandler) Update(e echo.Context) error {
	username := e.Param("username")
	c := new(model.User)
	if err := e.Bind(c); err != nil {
		return err
	}
	err := handler.service.Update(e.Request().Context(), username, c.Password, c.Admin)
	if err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	return e.JSON(http.StatusOK, c)
}
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
