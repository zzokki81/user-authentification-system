package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthChecker struct {
	interactor healthInteractor
}

type healthInteractor interface {
	CheckStoreConnection() error
}

func NewHealthChecker(interactor healthInteractor) *HealthChecker {
	return &HealthChecker{interactor}
}
func (hc *HealthChecker) HealthCheck(ctx echo.Context) error {
	if err := hc.interactor.CheckStoreConnection(); err != nil {
		return err
	}
	return ctx.String(http.StatusOK, "Welcome to initial web server")
}
