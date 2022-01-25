package handlers

import (
	"EgMeln/CRUDentity/internal/model"
	"EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ParkingLotHandler struct {
	Service *service.ParkingService
}

func NewServiceParkingLot(srv *service.ParkingService) ParkingLotHandler {
	return ParkingLotHandler{Service: srv}
}

func (handler *ParkingLotHandler) AddParkingLot(e echo.Context) (err error) {
	c := new(model.ParkingLot)
	if err = e.Bind(c); err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	err = handler.Service.AddParkingLot(e.Request().Context(), c)
	if err != nil {
		return e.JSON(http.StatusBadRequest, c)
	}
	return e.JSON(http.StatusOK, c)
}

func (handler *ParkingLotHandler) GetAllParkingLots(e echo.Context) error {
	parkingLots, err := handler.Service.GetAllParkingLots(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusBadRequest, parkingLots)
	}
	return e.JSON(http.StatusOK, parkingLots)
}

func (handler *ParkingLotHandler) GetParkingLotByNum(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	var parkingLot *model.ParkingLot
	parkingLot, err = handler.Service.GetParkingLotByNum(e.Request().Context(), num)
	if err != nil {
		return e.JSON(http.StatusBadRequest, parkingLot)
	}
	return e.JSON(http.StatusOK, parkingLot)
}

func (handler *ParkingLotHandler) UpdateParkingLot(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	c := new(model.ParkingLot)
	if err = e.Bind(c); err != nil {
		return err
	}
	car := model.ParkingLot{
		InParking: c.InParking,
		Remark:    c.Remark,
	}
	err = handler.Service.UpdateParkingLot(e.Request().Context(), num, car.InParking, car.Remark)
	if err != nil {
		return e.JSON(http.StatusBadRequest, car)
	}
	return e.JSON(http.StatusOK, car)
}

func (handler *ParkingLotHandler) DeleteParkingLot(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	err = handler.Service.DeleteParkingLot(e.Request().Context(), num)
	if err != nil {
		return e.JSON(http.StatusBadRequest, e)
	}
	return e.JSON(http.StatusOK, e)
}
