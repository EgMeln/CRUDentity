package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type ParkingService struct {
	conn repository.ParkingLots
}

func NewParkingLotServicePostgres(rep *repository.Postgres) *ParkingService {
	return &ParkingService{conn: rep}
}
func NewParkingLotServiceMongo(rep *repository.Mongo) *ParkingService {
	return &ParkingService{conn: rep}
}

func (srv *ParkingService) Add(e context.Context, lot *model.ParkingLot) error {
	return srv.conn.AddParkingLot(e, lot)
}
func (srv *ParkingService) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	return srv.conn.GetAllParkingLot(e)
}
func (srv *ParkingService) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	return srv.conn.GetByNumParkingLot(e, num)
}
func (srv *ParkingService) Update(e context.Context, num int, inParking bool, remark string) error {
	return srv.conn.UpdateParkingLot(e, num, inParking, remark)
}
func (srv *ParkingService) Delete(e context.Context, num int) error {
	return srv.conn.DeleteParkingLot(e, num)
}
