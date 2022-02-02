package main

import (
	"context"
	"fmt"
	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/handlers"
	"github.com/EgMeln/CRUDentity/internal/model"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln("Config error: ", cfg)
	}
	access := service.NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := service.NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	var parkingService *service.ParkingService
	var userService *service.UserService
	var authenticationService *service.AuthenticationService
	switch cfg.DB {
	case "postgres":
		cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, cfg.Host, cfg.PortPostgres, cfg.DBNamePostgres)
		log.Printf("DB URL: %s", cfg.DBURL)
		pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
		if err != nil {
			log.Fatalf("Error connection to DB: %v", err)
		}
		defer pool.Close()
		parkingService = service.NewParkingLotServicePostgres(&repository.Postgres{Pool: pool})
		userService = service.NewUserServicePostgres(&repository.Postgres{Pool: pool})
		authenticationService = service.NewAuthenticationServicePostgres(&repository.Postgres{Pool: pool}, access, refresh, cfg.HashSalt)
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
		parkingService = service.NewParkingLotServiceMongo(&repository.Mongo{CollectionParkingLot: db.Collection("egormelnikov")})
		userService = service.NewUserServiceMongo(&repository.Mongo{CollectionUsers: db.Collection("users")})
		authenticationService = service.NewAuthenticationServiceMongo(&repository.Mongo{CollectionUsers: db.Collection("users")}, access, refresh, cfg.HashSalt)
	}

	parkingHandler := handlers.NewServiceParkingLot(parkingService)
	userHandler := handlers.NewServiceUser(userService)
	authenticationHandler := handlers.NewServiceAuthentication(authenticationService)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/auth/sign-in", authenticationHandler.SignIn)
	e.POST("/auth/sign-up", authenticationHandler.SignUp)
	admin := e.Group("/admin")
	configuration := middleware.JWTConfig{Claims: &model.Claim{}, SigningKey: []byte(cfg.AccessToken)}
	admin.Use(middleware.JWTWithConfig(configuration))
	admin.Use(service.CheckAccess)
	admin.Use(service.TokenRefresh(access, refresh))

	admin.POST("/park", parkingHandler.Add)
	admin.PUT("/park/:num", parkingHandler.Update)
	admin.DELETE("/park/:num", parkingHandler.Delete)
	admin.GET("/users", userHandler.GetAll)
	admin.GET("/users/:username", userHandler.Get)
	admin.PUT("/users/:username", userHandler.Update)
	admin.DELETE("/users/:username", userHandler.Delete)

	user := e.Group("/user")
	user.Use(middleware.JWTWithConfig(configuration))
	user.Use(service.TokenRefresh(access, refresh))
	user.GET("/park", parkingHandler.GetAll)
	user.GET("/park/:num", parkingHandler.GetByNum)
	e.Logger.Fatal(e.Start(":8080"))
}
