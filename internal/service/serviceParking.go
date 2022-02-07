package service

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
)

// ParkingService struct for rep
type ParkingService struct {
	conn repository.ParkingLots
}

// NewParkingLotServicePostgres used for setting postgres services
func NewParkingLotServicePostgres(rep *repository.PostgresParking) *ParkingService {
	return &ParkingService{conn: rep}
}

// NewParkingLotServiceMongo used for setting mongo services
func NewParkingLotServiceMongo(rep *repository.MongoParking) *ParkingService {
	return &ParkingService{conn: rep}
}

// Add record about parking lot
func (srv *ParkingService) Add(e context.Context, lot *model.ParkingLot) error {
	return srv.conn.Add(e, lot)
}

// GetAll getting all parking lots
func (srv *ParkingService) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	return srv.conn.GetAll(e)
}

// GetByNum getting parking lot by num
func (srv *ParkingService) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	return srv.conn.GetByNum(e, num)
}

// Update updating parking lot
func (srv *ParkingService) Update(e context.Context, num int, inParking bool, remark string) error {
	return srv.conn.Update(e, num, inParking, remark)
}

// Delete deleting parking lot
func (srv *ParkingService) Delete(e context.Context, num int) error {
	return srv.conn.Delete(e, num)
}
