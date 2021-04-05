package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/stretchr/testify/require"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zzokki81/uas/interactor/mock"
	m "github.com/zzokki81/uas/middleware"
	"github.com/zzokki81/uas/model"
)

const (
	hashkey  = "secret"
	blockkey = "a-lot-secret-key"
)

var store = sessions.NewCookieStore([]byte(hashkey), []byte(blockkey))

func TestDo(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Expected response and actual response", func(t *testing.T) {
		interactor := &mock.Invitation{}
		userLoader := m.UserLoader{}
		userID := uint(1)
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/v1/users/1/invitation", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)
		c.Set("_session_store", store)
		sess, err := session.Get("session", c)
		require.NoError(err)
		sess.Values["userID"] = userID
		assert.NoError(sess.Save(c.Request(), c.Response()))

		h := func(c echo.Context) error {
			user := c.Get("user")
			return c.JSON(http.StatusOK, user)
		}

		expectedUser := &model.User{
			ID:        userID,
			Email:     "alex@edwards.com",
			Name:      "Alex",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		interactor.On("FindUserByID", userID).Return(expectedUser, nil)

		userLoader.Interactor = interactor

		mw := userLoader.Do(h)
		assert.NoError(mw(c))
		user := c.Get("user")

		assert.Equal(expectedUser, user)
	})
	t.Run("Error finding user with given ID", func(t *testing.T) {
		interactor := &mock.Invitation{}
		userLoader := m.UserLoader{}
		userID := uint(7)
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/v1/users/1/invitation", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)
		c.Set("_session_store", store)
		sess, err := session.Get("session", c)
		require.NoError(err)
		sess.Values["userID"] = userID
		assert.NoError(sess.Save(c.Request(), c.Response()))

		h := func(c echo.Context) error {
			return errors.New("cannot find user with given id")
		}

		interactor.On("FindUserByID", userID).Return(nil, errors.New("User with given ID can not be found"))

		userLoader.Interactor = interactor

		mw := userLoader.Do(h)
		assert.Error(mw(c))

		expectedError := errors.New("cannot find user with given id")

		assert.Equal(expectedError, h(c))
	})
	t.Run("No userID in session values", func(t *testing.T) {
		userLoader := m.UserLoader{}

		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/v1/users/1/invitation", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)
		c.Set("_session_store", store)
		sess, err := session.Get("session", c)
		require.NoError(err)
		assert.NoError(sess.Save(c.Request(), c.Response()))

		h := func(c echo.Context) error {
			user := c.Get("user")
			return c.JSON(http.StatusOK, user)
		}

		mw := userLoader.Do(h)
		assert.Error(mw(c))
	})

}
