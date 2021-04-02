package middleware

import (
	"errors"
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/zzokki81/uas/model"
)

type userInteractor interface {
	FindUserByID(id uint) (*model.User, error)
}

type UserLoader struct {
	Interactor userInteractor
}

func (ul UserLoader) Do(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		userID, ok := sess.Values["userID"].(uint)
		if !ok {
			return errors.New("no value for wanted key - userID")
		}

		user, err := ul.Interactor.FindUserByID(userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "cannot find user with given id")
		}
		c.Set("user", user)
		return next(c)
	}
}
