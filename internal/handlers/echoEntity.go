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
func (handler *ParkingLotHandler) Create(e echo.Context) (err error) {
	c := new(model.ParkingLot)
	if err = e.Bind(c); err != nil {
		return err
	}
	handler.Service.CreateRecord(c)
	return e.JSON(http.StatusOK, c)
}
func (handler *ParkingLotHandler) ReadAll(e echo.Context) error {
	var parkingLots []*model.ParkingLot
	parkingLots = handler.Service.ReadAllRecord()
	return e.JSON(http.StatusOK, parkingLots)
}
func (handler *ParkingLotHandler) ReadById(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	return e.JSON(http.StatusOK, handler.Service.ReadRecordByNum(num))
}
func (handler *ParkingLotHandler) UpdateRecord(e echo.Context) error {
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
	handler.Service.UpdateRecord(num, car.InParking, car.Remark)
	return e.JSON(http.StatusOK, e)
}
func (handler *ParkingLotHandler) DeleteRecord(e echo.Context) error {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	handler.Service.DeleteRecord(num)
	return e.JSON(http.StatusOK, e)
}

//func Create(e echo.Context) (err error) {
//	c := new(model.ParkingLot)
//	if err = e.Bind(c); err != nil {
//		return err
//	}
//	car := model.ParkingLot{
//		Num:       c.Num,
//		InParking: c.InParking,
//		Remark:    c.Remark,
//	}
//	repository.CreateRecord(car.Num, car.InParking, car.Remark)
//	return e.JSON(http.StatusOK, c)
//}
//
//func ReadAll(e echo.Context) error {
//	return e.String(http.StatusOK, repository.ReadAllRecords())
//}
//
//func ReadById(e echo.Context) (err error) {
//	num, err := strconv.Atoi(e.Param("num"))
//	if err != nil {
//		return err
//	}
//	return e.String(http.StatusOK, repository.ReadRecordByNum(num))
//}
//
//func UpdateRecord(e echo.Context) (err error) {
//	num, err := strconv.Atoi(e.Param("num"))
//	if err != nil {
//		return err
//	}
//	c := new(model.ParkingLot)
//	if err = e.Bind(c); err != nil {
//		return err
//	}
//	car := model.ParkingLot{
//		InParking: c.InParking,
//		Remark:    c.Remark,
//	}
//	repository.UpdateRecord(num, car.InParking, car.Remark)
//	return e.JSON(http.StatusOK, e)
//}
//
//func Delete(e echo.Context) (err error) {
//	num, err := strconv.Atoi(e.Param("num"))
//	if err != nil {
//		return err
//	}
//	repository.DeleteRecord(num)
//	return e.JSON(http.StatusOK, e)
//}
