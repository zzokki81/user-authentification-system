package interactor

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zzokki81/uas/handler/dto"
	"github.com/zzokki81/uas/model"
)

type Invitation struct {
	store Store
}

func NewInvitation(store Store) *Invitation {
	return &Invitation{
		store: store,
	}
}

func (i *Invitation) Create(request *dto.CreateInvitationRequest, user *model.User) (*model.Invitation, error) {
	invitation := &model.Invitation{
		Email:     request.Email,
		InviterID: user.ID,
	}

	tx, err := i.store.Transaction()
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Transaction did not open")
	}

	if err := tx.CreateInvitation(invitation); err != nil {
		tx.Rollback()
		return nil, err
	}

	notification := model.NewInvitationNotification(invitation.Email)
	if err := tx.CreateNotification(notification); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return invitation, nil
}

func (i *Invitation) FindByInviter(inviterID int) ([]*model.Invitation, error) {
	return i.store.FindInvitationByInviter(inviterID)
}

func (i *Invitation) FindUserByID(id uint) (*model.User, error) {
	return i.store.FindUserByID(id)
}

func (i *Invitation) FindAll() ([]*model.Invitation, error) {
	return i.store.FindAllInvitations()
}
