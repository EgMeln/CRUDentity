package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	postgresDB  *pgxpool.Pool
	dbClient    *mongo.Client
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	postgresPool, postgresResource := testPostgres()
	mongoPool, mongoResource := testMongo()
	redisPool, redisResource := testRedis()

	code := m.Run()

	if err := postgresPool.Purge(postgresResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	if err := mongoPool.Purge(mongoResource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
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
	return pool, resource
}
