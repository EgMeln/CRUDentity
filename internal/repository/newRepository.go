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
	CollectionUsers      *mongo.Collection
}

type ParkingLots interface {
	AddParkingLot(e context.Context, lot *model.ParkingLot) error
	GetAllParkingLot(e context.Context) ([]*model.ParkingLot, error)
	GetByNumParkingLot(e context.Context, num int) (*model.ParkingLot, error)
	UpdateParkingLot(e context.Context, num int, inParking bool, remark string) error
	DeleteParkingLot(e context.Context, num int) error
}
type Users interface {
	AddUser(e context.Context, lot *model.User) error
	GetAllUser(e context.Context) ([]*model.User, error)
	GetUser(e context.Context, username string) (*model.User, error)
	UpdateUser(e context.Context, username string, password string, admin bool) error
	DeleteUser(e context.Context, username string) error
}
