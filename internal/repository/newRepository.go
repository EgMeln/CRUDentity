package repository

import (
	"EgMeln/CRUDentity/internal/model"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

type Mongo struct {
	CollectionParkingLot *mongo.Collection
}

type ParkingLots interface {
	AddParkingLot(e context.Context, lot *model.ParkingLot) error
	GetAllParkingLots(e context.Context) ([]*model.ParkingLot, error)
	GetParkingLotByNum(e context.Context, num int) (*model.ParkingLot, error)
	UpdateParkingLot(e context.Context, num int, inParking bool, remark string) error
	DeleteParkingLot(e context.Context, num int) error
}
