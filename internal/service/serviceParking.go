package service

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

type ParkingService struct {
	conn repository.ParkingLots
}

func NewServicePostgres(rep *repository.Postgres) *ParkingService {
	return &ParkingService{conn: rep}
}
func NewServiceMongo(rep *repository.Mongo) *ParkingService {
	return &ParkingService{conn: rep}
}

//func NewService(rep interface{}) *ParkingService {
//	var postgres *repository.Postgres
//	var mongo *repository.Mongo
//	switch reflect.TypeOf(rep) {
//	case reflect.TypeOf(postgres):
//		return &ParkingService{conn: rep.(*repository.Postgres)}
//	case reflect.TypeOf(mongo):
//		return &ParkingService{conn: rep.(*repository.Mongo)}
//	}
//	return nil
//}

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
