package main

import (
	"EgMeln/CRUDentity/internal/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", handlers.Create)
	e.GET("/park", handlers.ReadAll)
	e.GET("/park/:num", handlers.ReadById)
	e.PUT("/change/:num", handlers.UpdateRecord)
	e.DELETE("/delete/:num", handlers.Delete)
	e.Logger.Fatal(e.Start(":1322"))
}
