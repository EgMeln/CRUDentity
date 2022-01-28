package repository

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
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
	Add(e context.Context, lot *model.ParkingLot) error
	GetAll(e context.Context) ([]*model.ParkingLot, error)
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Update(e context.Context, num int, inParking bool, remark string) error
	Delete(e context.Context, num int) error
}
