package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type ParkingService struct {
	conn repository.ParkingLots
}

func NewParkingLotServicePostgres(rep *repository.PostgresParking) *ParkingService {
	return &ParkingService{conn: rep}
}
func NewParkingLotServiceMongo(rep *repository.MongoParking) *ParkingService {
	return &ParkingService{conn: rep}
}

func (srv *ParkingService) Add(e context.Context, lot *model.ParkingLot) error {
	return srv.conn.Add(e, lot)
}
func (srv *ParkingService) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	return srv.conn.GetAll(e)
}
func (srv *ParkingService) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	return srv.conn.GetByNum(e, num)
}
func (srv *ParkingService) Update(e context.Context, num int, inParking bool, remark string) error {
	return srv.conn.Update(e, num, inParking, remark)
}
func (srv *ParkingService) Delete(e context.Context, num int) error {
	return srv.conn.Delete(e, num)
}
