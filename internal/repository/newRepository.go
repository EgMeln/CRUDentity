// Package repository contains code for handling different types of databases
package repository

import (
	"context"

	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostgresUser struct for user pool
type PostgresUser struct {
	PoolUser *pgxpool.Pool
}

// PostgresParking struct for parking pool
type PostgresParking struct {
	PoolParking *pgxpool.Pool
}

// PostgresToken struct for token pool
type PostgresToken struct {
	PoolToken *pgxpool.Pool
}

// MongoParking struct for parking pool
type MongoParking struct {
	CollectionParkingLot *mongo.Collection
}

// MongoUser struct for parking pool
type MongoUser struct {
	CollectionUsers *mongo.Collection
}

// MongoToken struct for parking pool
type MongoToken struct {
	CollectionTokens *mongo.Collection
}

// ParkingLots used for structuring, function for working with parking lots
type ParkingLots interface {
	Add(e context.Context, lot *model.ParkingLot) error
	GetAll(e context.Context) ([]*model.ParkingLot, error)
	GetByNum(e context.Context, num int) (*model.ParkingLot, error)
	Update(e context.Context, num int, inParking bool, remark string) error
	Delete(e context.Context, num int) error
}

// Users used for structuring, function for working with users
type Users interface {
	Add(e context.Context, lot *model.User) error
	GetAll(e context.Context) ([]*model.User, error)
	Get(e context.Context, username string) (*model.User, error)
	Update(e context.Context, user *model.User) error
	Delete(e context.Context, username string) error
}

// Authentication used for structuring, function for working with authentication
type Authentication interface {
	SignIn(e context.Context, user *model.User) (string, string, error)
}

// Tokens used for structuring, function for working with tokens
type Tokens interface {
	Add(e context.Context, token *model.Token) error
	Get(e context.Context, username string) (string, error)
	Delete(e context.Context, username string) error
}
