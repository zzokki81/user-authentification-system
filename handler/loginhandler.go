package handler

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/zzokki81/uas/model"
)

type LoginHandler struct {
	domain string
}

func NewLoginHandler(domain string) *LoginHandler {
	return &LoginHandler{domain}
}

func (lh *LoginHandler) Login(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Domain:   lh.domain,
		Path:     "/",
		MaxAge:   3600 * 8,
		HttpOnly: true,
	}
	user := &model.User{ID: 1, Email: "test@test.com", Name: "Ime"}
	sess.Values["userID"] = user.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
