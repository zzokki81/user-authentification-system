package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zzokki81/uas/handler"
)

const (
	hashkey        = "secret"
	blockkey       = "a-lot-secret-key"
	invalidHashKey = "sss"
	domain         = "localhost"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	t.Run("Login successfull", func(t *testing.T) {
		e := echo.New()
		handler := handler.NewLoginHandler(domain)
		req := httptest.NewRequest(http.MethodGet, "/v1/login", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()
		store := sessions.NewCookieStore([]byte(hashkey), []byte(blockkey))
		c := e.NewContext(req, rr)
		c.Set("_session_store", store)
		assert.NoError(handler.Login(c))
		var s = securecookie.New([]byte(hashkey), []byte(blockkey))
		cookie := rr.Result().Cookies()[0]

		values := make(map[interface{}]interface{})
		err := s.Decode(cookie.Name, cookie.Value, &values)
		assert.NoError(err)
		userID := values["userID"]
		assert.Equal(userID, uint(1))
	})
	t.Run("Invalid hash keys", func(t *testing.T) {
		e := echo.New()
		handler := handler.NewLoginHandler(domain)
		req := httptest.NewRequest(http.MethodGet, "/v1/login", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()
		store := sessions.NewCookieStore([]byte(invalidHashKey), []byte(invalidHashKey))
		c := e.NewContext(req, rr)
		c.Set("_session_store", store)
		assert.Error(handler.Login(c))
	})
}
