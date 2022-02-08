package handlers

import (
	"net/http"
	"strconv"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/request"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

// ParkingLotHandler struct that contain repository linc
type ParkingLotHandler struct {
	service *service.ParkingService
}

// NewServiceParkingLot add new authentication handler
func NewServiceParkingLot(srv *service.ParkingService) ParkingLotHandler {
	return ParkingLotHandler{service: srv}
}

// Add record about parking lot
func (handler *ParkingLotHandler) Add(e echo.Context) (err error) { //nolint:dupl //Different business logic
	c := new(request.ParkingLotCreate)
	if err = e.Bind(c); err != nil {
		log.WithField("Error", err).Warn("Bind fail")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Add(e.Request().Context(), &model.ParkingLot{Num: c.Num, InParking: c.InParking, Remark: c.Remark})
	if err != nil {
		log.WithField("Error", err).Warn("Add parking lot error")
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// GetAll getting all parking lots
func (handler *ParkingLotHandler) GetAll(e echo.Context) error {
	parkingLots, err := handler.service.GetAll(e.Request().Context())
	if err != nil {
		log.WithField("Error", err).Warn("Get all parking lots error")
		return echo.NewHTTPError(http.StatusInternalServerError, parkingLots)
	}
	return e.JSON(http.StatusOK, parkingLots)
}

// GetByNum getting parking lot by num
func (handler *ParkingLotHandler) GetByNum(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		log.WithField("Error", err).Warn("Num conv fail")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	var parkingLot *model.ParkingLot
	parkingLot, err = handler.service.GetByNum(e.Request().Context(), num)
	if err != nil {
		log.WithField("Error", err).Warn("Get by num parking lot error")
		return echo.NewHTTPError(http.StatusInternalServerError, parkingLot)
	}
	return e.JSON(http.StatusOK, parkingLot)
}

// Update updating parking lot
func (handler *ParkingLotHandler) Update(e echo.Context) error {
	c := new(request.ParkingLotUpdate)
	if err := e.Bind(c); err != nil {
		log.WithField("Error", err).Warn("Bind fail")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		log.WithField("Error", err).Warn("Num conv fail")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Update(e.Request().Context(), num, c.InParking, c.Remark)
	if err != nil {
		log.WithField("Error", err).Warn("Update parking lot error")
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting parking lot
func (handler *ParkingLotHandler) Delete(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		log.WithField("Error", err).Warn("Num conv fail")
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Delete(e.Request().Context(), num)
	if err != nil {
		log.WithField("Error", err).Warn("Delete parking lot error")
		return echo.NewHTTPError(http.StatusInternalServerError, e)
	}
	return e.JSON(http.StatusOK, e)
}
