package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/handlers"
	"github.com/EgMeln/CRUDentity/internal/middlewares"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/EgMeln/CRUDentity/internal/validation"
	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	initLog()

	cfg, err := config.New()
	if err != nil {
		log.Warnf("Config error %v", err)
	}
	access := service.NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := service.NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	var parkingService *service.ParkingService
	var userService *service.UserService
	var authenticationService *service.AuthenticationService
	switch cfg.DB {
	case "postgres":
		cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, cfg.Host, cfg.PortPostgres, cfg.DBNamePostgres)
		log.Infof("DB URL: %s", cfg.DBURL)
		pool := connectPostgres(cfg.DBURL)
		parkingService = service.NewParkingLotServicePostgres(&repository.PostgresParking{PoolParking: pool})
		userService = service.NewUserServicePostgres(&repository.PostgresUser{PoolUser: pool})
		authenticationService = service.NewAuthServicePostgres(&repository.PostgresToken{PoolToken: pool}, access, refresh, cfg.HashSalt)
	case "mongodb":
		cfg.DBURL = fmt.Sprintf("%s://%s:%d", cfg.DB, cfg.HostMongo, cfg.PortMongo)
		log.Infof("DB URL: %s", cfg.DBURL)
		client, db := connectMongo(cfg.DBURL, cfg.DBNameMongo)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				log.Warnf("Error connection to DB %v", err)
			}
		}()
		parkingService = service.NewParkingLotServiceMongo(&repository.MongoParking{CollectionParkingLot: db.Collection("egormelnikov")})
		userService = service.NewUserServiceMongo(&repository.MongoUser{CollectionUsers: db.Collection("users")})
		repMongoTokens := &repository.MongoToken{CollectionTokens: db.Collection("tokens")}
		authenticationService = service.NewAuthServiceMongo(repMongoTokens, access, refresh, cfg.HashSalt)
	}

	parkingHandler := handlers.NewServiceParkingLot(parkingService)
	userHandler := handlers.NewServiceUser(userService, authenticationService)

	runEcho(&parkingHandler, &userHandler, cfg)
}
func connectPostgres(URL string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), URL)
	if err != nil {
		log.Warnf("Error connection to DB %v", err)
	}
	defer pool.Close()
	return pool
}
func connectMongo(URL, DBName string) (*mongo.Client, *mongo.Database) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(URL))
	if err != nil {
		log.Warnf("Error connection to DB %v", err)
	}
	db := client.Database(DBName)
	return client, db
}
func runEcho(parkingHandler *handlers.ParkingLotHandler, userHandler *handlers.UserHandler, cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/auth/sign-in", userHandler.SignIn)
	e.POST("/auth/sign-up", userHandler.Add)

	admin := e.Group("/admin")
	configuration := middleware.JWTConfig{Claims: &config.Claim{}, SigningKey: []byte(cfg.AccessToken)}
	admin.Use(middleware.JWTWithConfig(configuration))
	admin.Use(middlewares.CheckAccess)

	admin.POST("/park", parkingHandler.Add)
	admin.PUT("/park/:num", parkingHandler.Update)
	admin.DELETE("/park/:num", parkingHandler.Delete)
	admin.GET("/users", userHandler.GetAll)
	admin.GET("/users/:username", userHandler.Get)
	admin.PUT("/users/:username", userHandler.Update)
	admin.DELETE("/users/:username", userHandler.Delete)
	user := e.Group("/user")
	user.Use(middleware.JWTWithConfig(configuration))

	user.POST("/refresh", userHandler.Refresh)
	user.GET("/park", parkingHandler.GetAll)
	user.GET("/park/:num", parkingHandler.GetByNum)
	e.Logger.Fatal(e.Start(":8080"))
	return e
}
func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
