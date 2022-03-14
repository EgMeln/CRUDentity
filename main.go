package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	_ "github.com/EgMeln/CRUDentity/docs"
	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/handlers"
	"github.com/EgMeln/CRUDentity/internal/middlewares"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/EgMeln/CRUDentity/internal/server"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/EgMeln/CRUDentity/internal/validation"
	"github.com/EgMeln/CRUDentity/protocol"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

// @title CRUD entity API
// @version 1.0
// @description CRUD entity API for Golang Project Parking.

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	initLog()

	cfg, err := config.New()
	if err != nil {
		log.Warnf("Config error %v", err)
	}
	access := service.NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := service.NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	redisCfg, err := config.NewRedis()
	if err != nil {
		log.Warnf("redis config error: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisCfg.Addr, Password: redisCfg.Password, DB: redisCfg.DB})

	var parkingService *service.ParkingService
	var userService *service.UserService
	var authenticationService *service.AuthenticationService
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch cfg.DB {
	case "postgres":
		cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, "localhost", cfg.PortPostgres, cfg.DBNamePostgres)
		log.Infof("DB URL: %s", cfg.DBURL)
		pool := connectPostgres(cfg.DBURL)
		parkingService = service.NewParkingLotServicePostgres(&repository.PostgresParking{PoolParking: pool}, repository.NewParkingLotCache(ctx, rdb))
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
		parkingService = service.NewParkingLotServiceMongo(&repository.MongoParking{CollectionParkingLot: db.Collection("egormelnikov")}, repository.NewParkingLotCache(ctx, rdb))
		userService = service.NewUserServiceMongo(&repository.MongoUser{CollectionUsers: db.Collection("users")})
		repMongoTokens := &repository.MongoToken{CollectionTokens: db.Collection("tokens")}
		authenticationService = service.NewAuthServiceMongo(repMongoTokens, access, refresh, cfg.HashSalt)
	}

	switch cfg.Server {
	case "echo":
		parkingHandler := handlers.NewServiceParkingLot(parkingService)
		userHandler := handlers.NewServiceUser(userService, authenticationService)
		fileHandler := handlers.ImageHandler{}
		runEcho(&parkingHandler, &userHandler, &fileHandler, cfg)
	case "grpc":
		imageStore := service.NewDiskImageStore("client")
		parkingServer := server.NewParkingServer(parkingService)
		userServer := server.NewUserServer(userService, authenticationService, imageStore)
		err = runGRPC(parkingServer, userServer, access, refresh)
	}
	log.Info("HTTP server terminated", err)
}
func connectPostgres(URL string) *pgxpool.Pool {
	pool, err := pgxpool.Connect(context.Background(), URL)
	if err != nil {
		log.Warnf("Error connection to DB %v", err)
	}
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
func runEcho(parkingHandler *handlers.ParkingLotHandler, userHandler *handlers.UserHandler, fileHandler *handlers.ImageHandler, cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.GET("/swagger/*any", swaggerFiles.WrapHandler)
	e.Validator = &validation.CustomValidator{Validator: validator.New()}

	e.POST("/auth/sign-in", userHandler.SignIn)
	e.POST("/auth/sign-up", userHandler.Add)
	admin := e.Group("/admin")
	configuration := middleware.JWTConfig{Claims: &config.Claim{}, SigningKey: []byte(cfg.AccessToken)}
	admin.Use(middleware.JWTWithConfig(configuration))
	admin.Use(middlewares.CheckAccess)

	admin.PUT("/park", parkingHandler.Update)
	admin.DELETE("/park/:num", parkingHandler.Delete)
	admin.GET("/users", userHandler.GetAll)
	admin.GET("/users/:username", userHandler.Get)
	admin.PUT("/users", userHandler.Update)
	admin.DELETE("/users/:username", userHandler.Delete)
	admin.POST("/park", parkingHandler.Add)
	user := e.Group("/user")
	user.Use(middleware.JWTWithConfig(configuration))

	e.POST("/refresh", userHandler.Refresh)
	user.GET("/park", parkingHandler.GetAll)

	user.GET("/park/:num", parkingHandler.GetByNum)
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})
	e.POST("/uploadImage", fileHandler.Upload)
	e.GET("/downloadImage", fileHandler.Download)

	e.Logger.Fatal(e.Start(":8080"))
	return e
}
func runGRPC(parkingServer protocol.ParkingServiceServer, userServer protocol.UserServiceServer, access, refresh *service.JWTService) error {
	listener, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	interceptors := server.NewAuthInterceptor(access, refresh)

	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptors.UnaryServerAuthInterceptor()),
		grpc.StreamInterceptor(interceptors.StreamServerAuthInterceptor()),
	}
	grpcServer := grpc.NewServer(serverOptions...)

	protocol.RegisterUserServiceServer(grpcServer, userServer)
	protocol.RegisterParkingServiceServer(grpcServer, parkingServer)

	log.Printf("server listening at %v", listener.Addr())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return grpcServer.Serve(listener)
}
func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
