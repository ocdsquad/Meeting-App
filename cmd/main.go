package main

import (
	"E-Meeting/configs"
	"E-Meeting/presenter/handler"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title E-Meeting API
// @version 1.0

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	//load config
	config := configs.LoadConfig()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	if err := handler.RoutingRestAPI(e, config); err != nil {
		e.Logger.Error(err)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.App.Port)))
}
