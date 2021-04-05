package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zzokki81/uas/handler"
	"github.com/zzokki81/uas/handler/dto"
	imock "github.com/zzokki81/uas/interactor/mock"
	"github.com/zzokki81/uas/model"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	user := &model.User{
		ID: 1,
	}

	t.Run("Invitation creation", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()
		requestBody := `{"email":"alex@edwards.com"}`
		req := httptest.NewRequest(http.MethodPost, "/v1/users/:user_id/invitation", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set("user", user)

		var ID uint = 1
		interactor.On("FindUserByID", ID).Return(&model.User{
			ID: ID,
		}, nil)

		createInvitationRequest := &dto.CreateInvitationRequest{
			Email: "alex@edwards.com",
		}

		invitation := &model.Invitation{
			ID:        1,
			Email:     "alex@edwards.com",
			InviterID: 1,
			CreatedAt: time.Time{},
		}

		interactor.On("Create", createInvitationRequest, user).Return(invitation, nil)

		err := handler.Create(c)
		assert.Nil(err)
		expectedResponse := dto.InvitationResponse{
			ID,
			"alex@edwards.com",
			1,
			time.Time{},
		}
		var response dto.InvitationResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)
		assert.Equal(expectedResponse, response)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()

		requestBody := `{email:"alex@edwards.com"}`
		req := httptest.NewRequest(http.MethodPost, "/invitation", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set("user", user)

		interactor.On("FindUserByID", user).Return(user, nil)

		err := handler.Create(c)
		expectedError := "code=400, message=Syntax error: offset=2, error=invalid character 'e' looking for beginning of object key string"
		assert.Contains(err.Error(), expectedError)
	})

	t.Run("Creation of invitation failed in database", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()

		requestBody := `{"email":"alex@edwards.com"}`
		req := httptest.NewRequest(http.MethodPost, "/invitation", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set("user", user)

		interactor.On("FindUserByID", user).Return(user, nil)

		email := "alex@edwards.com"

		interactor.On("Create", &dto.CreateInvitationRequest{
			Email: email,
		}, user).
			Return(nil, errors.New("Error creating invitation"))

		err := handler.Create(c)
		expectedResponse := "Error creating invitation"
		assert.Equal(expectedResponse, err.Error())
	})

	t.Run("Validation failed", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()

		requestBody := `{"email":"alex"}`
		req := httptest.NewRequest(http.MethodPost, "/invitation", strings.NewReader(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set("user", user)

		interactor.On("FindUserByID", user).Return(user, nil)
		email := "alex"

		interactor.On("Create", &dto.CreateInvitationRequest{
			Email: email,
		}, user).
			Return(nil, errors.New("Error creating invitation"))

		handler.Interactor = interactor

		err := handler.Create(c)
		expectedError := "code=400, message=Please enter valid email"
		assert.Contains(err.Error(), expectedError)
	})

	t.Run("Cannot retreive user from context", func(t *testing.T) {
		handler := handler.NewInvitationHandler(nil)

		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/invitation", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		err := handler.Create(c)
		expectedError := "code=500, message=error retrieving user from context, internal=<nil>"
		assert.Equal(expectedError, err.Error())
	})

}

func TestFindInvitationByInviter(t *testing.T) {
	assert := assert.New(t)

	require := require.New(t)
	t.Run("Successfully returned invitations for given inviter id", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()

		expectedInvitations := []*dto.InvitationResponse{
			{
				ID:        1,
				Email:     "alex@edwards.com",
				InviterID: 1,
				CreatedAt: time.Time{},
			},
			{
				ID:        2,
				Email:     "jon@kalhun.com",
				InviterID: 1,
				CreatedAt: time.Time{},
			},
		}

		q := make(url.Values)
		q.Set("inviter_id", "1")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)
		c.QueryParam("inviter_id")

		invitation := []*model.Invitation{
			{
				ID:        1,
				Email:     "alex@edwards.com",
				InviterID: 1,
				CreatedAt: time.Time{},
			},
			{
				ID:        2,
				Email:     "jon@kalhun.com",
				InviterID: 1,
				CreatedAt: time.Time{},
			},
		}

		interactor.On("FindByInviter", 1).Return(invitation, nil)

		err := handler.FindInvitationByInviter(c)
		assert.Nil(err)
		assert.Equal(http.StatusOK, rr.Code)

		var response []*dto.InvitationResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)

		assert.Equal(expectedInvitations, response)

	})

	t.Run("Getting invitations failed in database", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		e := echo.New()

		q := make(url.Values)
		q.Set("inviter_id", "1")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		interactor.On("FindByInviter", 1).Return(nil, errors.New("Error retrieving invitations"))

		handler.Interactor = interactor

		err := handler.FindInvitationByInviter(c)
		expectedResponse := "Error retrieving invitations"
		assert.Equal(expectedResponse, err.Error())
	})
	t.Run("error converting inviter id to integer", func(t *testing.T) {
		interactor := &imock.Invitation{}
		handler := handler.NewInvitationHandler(interactor)

		expectedError := "strconv.Atoi: parsing \"ww\": invalid syntax"

		e := echo.New()

		q := make(url.Values)
		q.Set("inviter_id", "ww")
		req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		handler.Interactor = interactor
		err := handler.FindInvitationByInviter(c)
		assert.NotNil(err)

		assert.Equal(expectedError, err.Error())

	})
}
