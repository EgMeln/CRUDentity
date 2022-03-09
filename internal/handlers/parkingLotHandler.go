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
// @Summary add parking lot
// @ID add-parkingLot
// @Security ApiKeyAuth
// @Produce json
// @Param request body request.ParkingLotCreate true "create parking lot"
// @Success 200 {object} request.ParkingLotCreate
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/park [post]
func (handler *ParkingLotHandler) Add(e echo.Context) (err error) { //nolint:dupl //Different business logic
	c := new(request.ParkingLotCreate)
	if err = e.Bind(c); err != nil {
		log.Warnf("Bind fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if ok := e.Validate(c); ok != nil {
		log.Warnf("Validation error: %v", ok)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Add(e.Request().Context(), &model.ParkingLot{Num: c.Num, InParking: c.InParking, Remark: c.Remark})
	if err != nil {
		log.Warnf("Add parking lot error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// GetAll getting all parking lots
// @Summary gets all parking lots
// @ID get-all-parkingLots
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {array} model.ParkingLot
// @Failure 500 {array} model.ParkingLot
// @Router /user/park [get]
func (handler *ParkingLotHandler) GetAll(e echo.Context) error {
	parkingLots, err := handler.service.GetAll(e.Request().Context())
	if err != nil {
		log.Warnf("Get all parking lots error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, parkingLots)
	}
	returnParkingLot := make([]*request.ParkingLotsReturn, len(parkingLots))
	for i := range parkingLots {
		returnParkingLot = append(returnParkingLot, &request.ParkingLotsReturn{Num: (parkingLots[i]).Num, InParking: (parkingLots[i]).InParking, Remark: (parkingLots[i]).Remark})
	}
	return e.JSON(http.StatusOK, returnParkingLot)
}

// GetByNum getting parking lot by num
// @Summary get parking lot by num
// @ID get-parkingLot-by-num
// @Security ApiKeyAuth
// @Produce json
// @Param num path string true "get parking lot"
// @Success 200 {object} model.ParkingLot
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /user/park/{num} [get]
func (handler *ParkingLotHandler) GetByNum(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		log.Warnf("Num conv fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	var parkingLot *model.ParkingLot
	parkingLot, err = handler.service.GetByNum(e.Request().Context(), num)
	if err != nil {
		log.Warnf("Get by num parking lot error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, parkingLot)
	}
	returnParkingLot := &request.ParkingLotsReturn{Num: parkingLot.Num, InParking: parkingLot.InParking, Remark: parkingLot.Remark}
	return e.JSON(http.StatusOK, returnParkingLot)
}

// Update updating parking lot
// @Summary update parking lot by num
// @ID update-parkingLot-by-num
// @Security ApiKeyAuth
// @Produce json
// @Param request body request.ParkingLotUpdate true "update parking lot"
// @Success 200 {object} request.ParkingLotUpdate
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/park [put]
func (handler *ParkingLotHandler) Update(e echo.Context) error {
	c := new(request.ParkingLotUpdate)
	if err := e.Bind(c); err != nil {
		log.Warnf("Bind fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := e.Validate(c); err != nil {
		log.Warnf("Validation error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err := handler.service.Update(e.Request().Context(), c.Num, c.InParking, c.Remark)
	if err != nil {
		log.Warnf("Update parking lot error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, c)
	}
	return e.JSON(http.StatusOK, c)
}

// Delete deleting parking lot
// @Summary delete parking lot by num
// @ID delete-parkingLot-by-num
// @Security ApiKeyAuth
// @Produce json
// @Param num path string true "delete parking lot"
// @Success 200 {string} echo.Context
// @Failure 400 {string} echo.NewHTTPError
// @Failure 500 {string} echo.NewHTTPError
// @Router /admin/park/{num} [delete]
func (handler *ParkingLotHandler) Delete(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		log.Warnf("Num conv fail %v", err)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	err = handler.service.Delete(e.Request().Context(), num)
	if err != nil {
		log.Warnf("Delete parking lot error %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, e)
	}
	return e.JSON(http.StatusOK, e)
}
