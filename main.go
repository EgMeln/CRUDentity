package main

import (
	"EgMeln/CRUDentity/internal/config"
	"EgMeln/CRUDentity/internal/handlers"
	"EgMeln/CRUDentity/internal/repository"
	"EgMeln/CRUDentity/internal/service"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalln("Config error: ", cfg)
	}
	//cfg.DBURL = fmt.Sprintf("%s://%s:%s@%s:%d/%s", cfg.DB, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	//log.Printf("DB URL: %s", cfg.DBURL)
	//pool, err := pgxpool.Connect(context.Background(), cfg.DBURL)
	//if err != nil {
	//	log.Fatalf("Error connection to DB: %v", err)
	//}
	//defer pool.Close()
	connRepository := repository.NewRepository(cfg)
	parkingService := service.New(connRepository)
	parkingHandler := handlers.NewServiceParkingLot(parkingService)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", parkingHandler.Create)
	e.GET("/park", parkingHandler.ReadAll)
	e.GET("/park/:num", parkingHandler.ReadById)
	e.PUT("/change/:num", parkingHandler.UpdateRecord)
	e.DELETE("/delete/:num", parkingHandler.DeleteRecord)
	e.Logger.Fatal(e.Start(":1323"))
}
