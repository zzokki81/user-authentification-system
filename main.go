package main

import (
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/zzokki81/uas/handler"
	"github.com/zzokki81/uas/interactor"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	healthInteractor := interactor.NewHealthCheck()
	healthCheckHandler := handler.NewHealthChecker(healthInteractor)
	e := echo.New()
	e.Use(echomiddleware.Logger())

	e.GET("/v1/health", healthCheckHandler.HealthCheck)
	e.Server.Addr = ":8080"
	graceful.DefaultLogger().Println("Application has successfully started at port: ", 8080)
	graceful.ListenAndServe(e.Server, 5*time.Second)
}
