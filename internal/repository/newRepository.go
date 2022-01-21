package repository

import (
	"EgMeln/CRUDentity/internal/config"
	"EgMeln/CRUDentity/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Postgres struct {
	pool *pgxpool.Pool
}

type ParkingLots interface {
	CreateRecord(lot *model.ParkingLot)
	ReadAllRecords() []*model.ParkingLot
	ReadRecordByNum(num int) *model.ParkingLot
	UpdateRecord(num int, inParking bool, remark string)
	DeleteRecord(num int)
}

func NewRepository(cfg *config.Config) *Postgres {
	cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
	if err != nil {
		log.Fatalf("Error connection to DB: %v", err)
	}
	return &Postgres{pool: pool}
}
