package repository

import (
	"context"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostgresUser struct {
	PoolUser *pgxpool.Pool
}

type PostgresParking struct {
	PoolParking *pgxpool.Pool
}
type PostgresToken struct {
	PoolToken *pgxpool.Pool
}
type MongoParking struct {
	CollectionParkingLot *mongo.Collection
}
type MongoUser struct {
	CollectionUsers *mongo.Collection
}
type MongoToken struct {
	CollectionTokens *mongo.Collection
}

type ParkingLots interface {
	Add(e context.Context, lot *model.ParkingLot) error
	GetAll(e context.Context) ([]*model.ParkingLot, error)
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Update(e context.Context, num int, inParking bool, remark string) error
	Delete(e context.Context, num int) error
}
type Users interface {
	Add(e context.Context, lot *model.User) error
	GetAll(e context.Context) ([]*model.User, error)
	Get(e context.Context, username string) (*model.User, error)
	Update(e context.Context, username string, password string, admin bool) error
	Delete(e context.Context, username string) error
}

type Authentication interface {
	SignUp(e context.Context, user *model.User) error
	SignIn(e context.Context, user *model.User) (string, string, error)
}
type Tokens interface {
	Add(e context.Context, token *model.Token) error
	Get(e context.Context, username string) (string, error)
	Delete(e context.Context, username string) error
}
