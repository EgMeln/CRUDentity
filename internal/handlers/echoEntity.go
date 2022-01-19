package handlers

import (
	"EgMeln/CRUDentity/internal/repository/postgreSQL"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ParkingLot struct {
	Num       int    `json:"num" form:"num" query:"num"`
	InParking bool   `json:"parking" form:"parking" query:"parking"`
	Remark    string `json:"remark" form:"remark" query:"remark"`
}

func Create(e echo.Context) (err error) {
	c := new(ParkingLot)
	if err = e.Bind(c); err != nil {
		return err
	}
	car := ParkingLot{
		Num:       c.Num,
		InParking: c.InParking,
		Remark:    c.Remark,
	}
	postgreSQL.CreateRecord(car.Num, car.InParking, car.Remark)
	return e.JSON(http.StatusOK, c)
}

func ReadAll(e echo.Context) error {
	return e.String(http.StatusOK, postgreSQL.ReadAllRecords())
}
func ReadById(e echo.Context) (err error) {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	return e.String(http.StatusOK, postgreSQL.ReadRecordByNum(num))
}
func UpdateRecord(e echo.Context) (err error) {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	c := new(ParkingLot)
	if err = e.Bind(c); err != nil {
		return err
	}
	car := ParkingLot{
		InParking: c.InParking,
		Remark:    c.Remark,
	}
	postgreSQL.UpdateRecord(num, car.InParking, car.Remark)
	return e.JSON(http.StatusOK, e)
}

func Delete(e echo.Context) (err error) {
	num, err := strconv.Atoi(e.Param("num"))
	if err != nil {
		return err
	}
	postgreSQL.DeleteRecord(num)
	return e.JSON(http.StatusOK, e)
}
