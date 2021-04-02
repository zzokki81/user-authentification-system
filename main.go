package main

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/zzokki81/uas/handler"
	"github.com/zzokki81/uas/interactor"
	"github.com/zzokki81/uas/middleware"
	"github.com/zzokki81/uas/store/postgres"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(err)
	}
	config := postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	store, err := postgres.Open(config)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	e := echo.New()
	e.Use(echomiddleware.Logger())
	userInteractor := interactor.NewUser(store)
	userLoader := middleware.UserLoader{
		Interactor: userInteractor,
	}

	healthInteractor := interactor.NewHealthCheck(store)
	healthCheckHandler := handler.NewHealthChecker(healthInteractor)

	invitationInteractor := interactor.NewInvitation(store)
	invitationHandler := handler.NewInvitationHandler(invitationInteractor)

	e.GET("/v1/health", healthCheckHandler.HealthCheck)

	e.GET("/invitations", invitationHandler.FindInvitationByInviter)
	e.POST("/v1/invitations", invitationHandler.Create, userLoader.Do)

	e.Server.Addr = ":8080"
	graceful.DefaultLogger().Println("Application has successfully started at port: ", 8080)
	graceful.ListenAndServe(e.Server, 5*time.Second)
}
