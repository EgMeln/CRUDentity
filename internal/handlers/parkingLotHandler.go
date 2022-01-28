package handlers

import (
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/request"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ParkingLotHandler struct {
	service *service.ParkingService
}

func NewServiceParkingLot(srv *service.ParkingService) ParkingLotHandler {
	return ParkingLotHandler{service: srv}
}

func (handler *ParkingLotHandler) Add(e echo.Context) (err error) {
	c := new(request.ParkingLotCreate)
	if err = e.Bind(c); err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	err = handler.service.Add(e.Request().Context(), &model.ParkingLot{Num: c.Num, InParking: c.InParking, Remark: c.Remark})
	if err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	return e.JSON(http.StatusOK, c)
}

func (handler *ParkingLotHandler) GetAll(e echo.Context) error {
	parkingLots, err := handler.service.GetAll(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusBadRequest, parkingLots)
	}
	return e.JSON(http.StatusOK, parkingLots)
}

func (handler *ParkingLotHandler) GetByNum(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	var parkingLot *model.ParkingLot
	parkingLot, err = handler.service.GetByNum(e.Request().Context(), num)
	if err != nil {
		return e.JSON(http.StatusBadRequest, parkingLot)
	}
	return e.JSON(http.StatusOK, parkingLot)
}

func (handler *ParkingLotHandler) Update(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	c := new(request.ParkingLotUpdate)
	if err = e.Bind(c); err != nil {
		return err
	}
	err = handler.service.Update(e.Request().Context(), num, c.InParking, c.Remark)
	if err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	return e.JSON(http.StatusOK, c)
}

func (handler *ParkingLotHandler) Delete(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	err = handler.service.Delete(e.Request().Context(), num)
	if err != nil {
		return e.JSON(http.StatusBadRequest, e)
	}
	return e.JSON(http.StatusOK, e)
}
