package service

import (
	"EgMeln/CRUDentity/internal/model"
	"EgMeln/CRUDentity/internal/repository"
)

type ParkingService struct {
	Conn *repository.Postgres
}

type ParkingLots interface {
	CreateRecord(lot *model.ParkingLot)
	ReadAllRecord() []*model.ParkingLot
	ReadRecordByNum(num int) *model.ParkingLot
	UpdateRecord(num int, inParking bool, remark string)
	DeleteRecord(num int)
}

func New(Rep *repository.Postgres) *ParkingService {
	return &ParkingService{Conn: Rep}
}
func (srv *ParkingService) CreateRecord(lot *model.ParkingLot) {
	srv.Conn.CreateRecord(lot)
}
func (srv *ParkingService) ReadAllRecord() []*model.ParkingLot {
	return srv.Conn.ReadAllRecords()
}
func (srv *ParkingService) ReadRecordByNum(num int) *model.ParkingLot {
	return srv.Conn.ReadRecordByNum(num)
}
func (srv *ParkingService) UpdateRecord(num int, inParking bool, remark string) {
	srv.Conn.UpdateRecord(num, inParking, remark)
}
func (srv *ParkingService) DeleteRecord(num int) {
	srv.Conn.DeleteRecord(num)
}
