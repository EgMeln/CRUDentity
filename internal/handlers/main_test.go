package handlers

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/EgMeln/CRUDentity/internal/config"
	"github.com/EgMeln/CRUDentity/internal/middlewares"
	"github.com/EgMeln/CRUDentity/internal/repository"
	"github.com/EgMeln/CRUDentity/internal/service"
	"github.com/EgMeln/CRUDentity/internal/validation"
	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	swaggerFiles "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	postgresDB   *pgxpool.Pool
	dbClient     *mongo.Client
	redisClient  *redis.Client
	accessToken  string
	refreshToken string
)

func TestMain(m *testing.M) {
	postgresPool, postgresResource := testPostgres()
	mongoPool, mongoResource := testMongo()
	redisPool, redisResource := testRedis()

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error %v", err)
	}
	access := service.NewJWTService([]byte(cfg.AccessToken), time.Duration(cfg.AccessTokenLifeTime)*time.Second)
	refresh := service.NewJWTService([]byte(cfg.RefreshToken), time.Duration(cfg.RefreshTokenLifeTime)*time.Second)

	hash, err := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("can't create hash password for user")
	}
	_, err = postgresDB.Exec(context.Background(), "INSERT INTO users (username,password,is_admin) VALUES ($1,$2,$3)", "handler", string(hash), true)
	if err != nil {
		log.Fatalf("can't insert admin to db")
	}
	ctx, cancel := context.WithCancel(context.Background())
	var parkingService *service.ParkingService
	var userService *service.UserService
	var authenticationService *service.AuthenticationService
	switch cfg.DB {
	case "postgres":
		parkingService = service.NewParkingLotServicePostgres(
			&repository.PostgresParking{PoolParking: postgresDB}, repository.NewParkingLotCache(ctx, redisClient))
		userService = service.NewUserServicePostgres(&repository.PostgresUser{PoolUser: postgresDB})
		authenticationService = service.NewAuthServicePostgres(
			&repository.PostgresToken{PoolToken: postgresDB}, access, refresh, cfg.HashSalt)
	case "mongodb":
		parkingService = service.NewParkingLotServiceMongo(
			&repository.MongoParking{CollectionParkingLot: dbClient.Database("egormelnikovdb").Collection("egormelnikov")}, repository.NewParkingLotCache(ctx, redisClient))
		userService = service.NewUserServiceMongo(&repository.MongoUser{CollectionUsers: dbClient.Database("egormelnikovdb").Collection("users")})
		repMongoTokens := &repository.MongoToken{CollectionTokens: dbClient.Database("egormelnikovdb").Collection("tokens")}
		authenticationService = service.NewAuthServiceMongo(repMongoTokens, access, refresh, cfg.HashSalt)
	}
	userHandler := NewServiceUser(userService, authenticationService)
	parkingHandler := NewServiceParkingLot(parkingService)
	fileHandler := ImageHandler{}
	runEcho(&parkingHandler, &userHandler, &fileHandler, cfg)

	code := m.Run()
	if err := postgresPool.Purge(postgresResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := mongoPool.Purge(mongoResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	cancel()
	if err := redisPool.Purge(redisResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := dbClient.Disconnect(context.TODO()); err != nil {
		log.Fatalf("mongo disconnection error %v", err)
	}
	os.Exit(code)
}
func testPostgres() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=54236305",
			"POSTGRES_USER=egormelnikov",
			"POSTGRES_DB=egormelnikov",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("8081/tcp")
	databaseURL := fmt.Sprintf("postgres://egormelnikov:54236305@%s/egormelnikov?sslmode=disable", hostAndPort)
	log.Println("Connecting to database on url: ", databaseURL)
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		postgresDB, err = pgxpool.Connect(context.Background(), databaseURL)
		if err != nil {
			return err
		}
		return postgresDB.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource
}
func testMongo() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	err = pool.Retry(func() error {
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return pool, resource
}

func testRedis() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	if err = pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		})

		return redisClient.Ping(context.Background()).Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	err = resource.Expire(60)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	return pool, resource
}
func runEcho(parkingHandler *ParkingLotHandler, userHandler *UserHandler, fileHandler *ImageHandler, cfg *config.Config) *echo.Echo {
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

	go func() {
		err := e.Start(":8081")
		if err != nil {
			log.Fatalf("error with starting an echo server: %v", err)
			return
		}
	}()
	time.Sleep(time.Second)
	return e
}
