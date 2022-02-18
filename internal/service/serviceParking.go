package service

import (
	"context"
	"fmt"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/labstack/gommon/log"
)

// ParkingService struct for rep
type ParkingService struct {
	conn  repository.ParkingLots
	redis repository.ParkingLotCacheRedis
}

// NewParkingLotServicePostgres used for setting postgres services
func NewParkingLotServicePostgres(rep *repository.PostgresParking, red *repository.ParkingLotCache) *ParkingService {
	return &ParkingService{conn: rep, redis: red}
}

// NewParkingLotServiceMongo used for setting mongo services
func NewParkingLotServiceMongo(rep *repository.MongoParking, red *repository.ParkingLotCache) *ParkingService {
	return &ParkingService{conn: rep, redis: red}
}

// Add record about parking lot
func (srv *ParkingService) Add(e context.Context, lot *model.ParkingLot) error {
	err := srv.conn.Add(e, lot)
	if err != nil {
		return fmt.Errorf("parking lot service %w", err)
	}
	err = srv.redis.Add(e, lot)
	if err != nil {
		return fmt.Errorf("cache parking lot service error %w", err)
	}
	return err
}

// GetAll getting all parking lots
func (srv *ParkingService) GetAll(e context.Context) ([]*model.ParkingLot, error) {
	return srv.conn.GetAll(e)
}

// GetByNum getting parking lot by num
func (srv *ParkingService) GetByNum(e context.Context, num int) (*model.ParkingLot, error) {
	lotCache, ok := srv.redis.GetByNum(e, num)
	if ok == nil {
		return lotCache, nil
	}
	log.Info("Cache is empty")
	lot, err := srv.conn.GetByNum(e, num)
	if err != nil {
		return nil, fmt.Errorf("parking lot service %w", err)
	}
	return lot, ok
}

// Update updating parking lot
func (srv *ParkingService) Update(e context.Context, num int, inParking bool, remark string) error {
	err := srv.conn.Update(e, num, inParking, remark)
	if err != nil {
		return fmt.Errorf("parking lot service %w", err)
	}
	err = srv.redis.Delete(e, num)
	if err != nil {
		return fmt.Errorf("cache parking lot service %w", err)
	}
	return err
}

// Delete deleting parking lot
func (srv *ParkingService) Delete(e context.Context, num int) error {
	err := srv.conn.Delete(e, num)
	if err != nil {
		return fmt.Errorf("parking lot service %w", err)
	}
	err = srv.redis.Delete(e, num)
	if err != nil {
		return fmt.Errorf("cache parking lot service %w", err)
	}
	return err
}
