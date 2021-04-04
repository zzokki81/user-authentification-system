package main

import (
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/zzokki81/uas/handler"
	"github.com/zzokki81/uas/interactor"
	"github.com/zzokki81/uas/middleware"
	"github.com/zzokki81/uas/service/oauth"
	"github.com/zzokki81/uas/store/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

	sessionAuthKey := os.Getenv("SESSION_AUTH_KEY")
	sessionEncryptionKey := os.Getenv("SESSION_ENCRYPTION_KEY")
	rawSessionAuthKey, err := base64.StdEncoding.DecodeString(sessionAuthKey)
	if err != nil {
		panic(err)
	}
	rawSessionEncryptionKey, err := base64.StdEncoding.DecodeString(sessionEncryptionKey)
	if err != nil {
		panic(err)
	}

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
	loginHandler := handler.NewLoginHandler(os.Getenv("DOMAIN"))

	googleOauthConfig := &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	google := oauth.NewGoogle(googleOauthConfig)
	googleHandler := handler.NewGoogleHandler(google)
	e.Use(session.Middleware(sessions.NewCookieStore(rawSessionAuthKey, rawSessionEncryptionKey)))

	e.GET("/v1/health", healthCheckHandler.HealthCheck)

	e.GET("/invitations", invitationHandler.FindInvitationByInviter)
	e.POST("/v1/invitations", invitationHandler.Create, userLoader.Do)

	e.GET("/v1/healthz", healthCheckHandler.HealthCheck)
	e.GET("/v1/login", loginHandler.Login)
	e.GET("/v1/google/login", googleHandler.HandleMain)
	e.GET("/v1/google/exchange", googleHandler.ExchangeAuthCodeForTOken)

	e.Server.Addr = ":8080"
	graceful.DefaultLogger().Println("Application has successfully started at port: ", 8080)
	graceful.ListenAndServe(e.Server, 5*time.Second)
}
