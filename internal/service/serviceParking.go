package service

import (
	"EgMeln/CRUDentity/internal/model"
	"EgMeln/CRUDentity/internal/repository"
	"context"
	"reflect"
)

type ParkingService struct {
	Conn repository.ParkingLots
}

func NewService(Rep interface{}) *ParkingService {
	var postgres *repository.Postgres
	var mongo *repository.Mongo
	switch reflect.TypeOf(Rep) {
	case reflect.TypeOf(postgres):
		return &ParkingService{Conn: Rep.(*repository.Postgres)}
	case reflect.TypeOf(mongo):
		return &ParkingService{Conn: Rep.(*repository.Mongo)}
	}
	return nil
}

func (srv *ParkingService) AddParkingLot(e context.Context, lot *model.ParkingLot) error {
	return srv.Conn.AddParkingLot(e, lot)
}
func (srv *ParkingService) GetAllParkingLots(e context.Context) ([]*model.ParkingLot, error) {
	return srv.Conn.GetAllParkingLots(e)
}
func (srv *ParkingService) GetParkingLotByNum(e context.Context, num int) (*model.ParkingLot, error) {
	return srv.Conn.GetParkingLotByNum(e, num)
}
func (srv *ParkingService) UpdateParkingLot(e context.Context, num int, inParking bool, remark string) error {
	return srv.Conn.UpdateParkingLot(e, num, inParking, remark)
}
func (srv *ParkingService) DeleteParkingLot(e context.Context, num int) error {
	return srv.Conn.DeleteParkingLot(e, num)
}
