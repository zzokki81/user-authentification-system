package interactor_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zzokki81/uas/handler/dto"
	"github.com/zzokki81/uas/interactor"
	"github.com/zzokki81/uas/model"
	smock "github.com/zzokki81/uas/store/mock"
)

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	t.Run("Invitation successfully saved in database", func(t *testing.T) {
		store := &smock.Store{}

		mockInvitation := &model.Invitation{
			Email:     "alex@edwards.com",
			InviterID: uint(1),
		}
		mockNotification := model.NewInvitationNotification(mockInvitation.Email)

		store.On("Transaction").Return(store, nil)
		store.On("CreateInvitation", mockInvitation).Return(nil)
		store.On("CreateNotification", mockNotification).Return(nil)
		store.On("Commit").Return(nil)

		createInvitationRequest := &dto.CreateInvitationRequest{
			Email: "alex@edwards.com",
		}
		user := &model.User{
			ID: uint(1),
		}

		i := interactor.NewInvitation(store)
		expectedInvitation, err := i.Create(createInvitationRequest, user)
		assert.NoError(err)

		assert.Equal(mockInvitation, expectedInvitation)
	})
	t.Run("Error creating invitation in database", func(t *testing.T) {
		store := &smock.Store{}

		mockInvitation := &model.Invitation{
			Email:     "alex@edwards.com",
			InviterID: uint(1),
		}

		mockNotification := &model.Notification{
			RecipientEmail: "alex@edwards.com",
			Type:           "Invitation",
		}
		storeError := errors.New("Error saving invitation in database")

		store.On("Transaction").Return(store, nil)
		store.On("CreateInvitation", mockInvitation).Return(storeError)
		store.On("CreateNotification", mockNotification).Return(nil)
		store.On("Rollback")

		l := interactor.NewInvitation(store)

		createInvitationRequest := &dto.CreateInvitationRequest{
			Email: "alex@edwards.com",
		}
		user := &model.User{
			ID: uint(1),
		}
		_, err := l.Create(createInvitationRequest, user)

		assert.Equal(err, storeError)
	})
	t.Run("Error opening transaction", func(t *testing.T) {
		store := &smock.Store{}

		storeError := "code=500, message=Transaction did not open"
		store.On("Transaction").Return(nil, errors.New("Transaction did not open"))

		l := interactor.NewInvitation(store)

		createInvitationRequest := &dto.CreateInvitationRequest{
			Email: "alex@edwards.com",
		}
		user := &model.User{
			ID: uint(1),
		}
		_, err := l.Create(createInvitationRequest, user)

		assert.Equal(storeError, err.Error())
	})
}
