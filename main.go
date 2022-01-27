package main

import (
	"EgMeln/CRUDentity/internal/config"
	"EgMeln/CRUDentity/internal/handlers"
	"EgMeln/CRUDentity/internal/repository"
	"EgMeln/CRUDentity/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln("Config error: ", cfg)
	}

	var parkingService *service.ParkingService
	switch cfg.DB {
	case "postgres":
		cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, cfg.Host, cfg.PortPostgres, cfg.DBNamePostgres)
		log.Printf("DB URL: %s", cfg.DBURL)
		pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
		if err != nil {
			log.Fatalf("Error connection to DB: %v", err)
		}
		defer pool.Close()
		parkingService = service.NewService(&repository.Postgres{Pool: pool})

	case "mongodb":
		cfg.DBURL = fmt.Sprintf("%s://%s:%d", cfg.DB, cfg.HostMongo, cfg.PortMongo)
		log.Printf("DB URL: %s", cfg.DBURL)
		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.DBURL))
		if err != nil {
			log.Fatalf("Error connection to DB: %v", err)
		}
		db := client.Database(cfg.DBNameMongo)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				log.Fatalf("Error connection to DB: %v", err)
			}
		}()
		parkingService = service.NewService(&repository.Mongo{CollectionParkingLot: db.Collection("egormelnikov")})
	}

	parkingHandler := handlers.NewServiceParkingLot(parkingService)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", parkingHandler.AddParkingLot)
	e.GET("/park", parkingHandler.GetAllParkingLots)
	e.GET("/park/:num", parkingHandler.GetParkingLotByNum)
	e.PUT("/change/:num", parkingHandler.UpdateParkingLot)
	e.DELETE("/delete/:num", parkingHandler.DeleteParkingLot)
	e.Logger.Fatal(e.Start(":8080"))
}
