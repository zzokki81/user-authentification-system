package handler

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"

	"github.com/zzokki81/uas/handler/dto"
	"github.com/zzokki81/uas/model"
)

type InvitationHandler struct {
	Interactor invitationInteractor
}

func NewInvitationHandler(invitationInteractor invitationInteractor) *InvitationHandler {
	return &InvitationHandler{invitationInteractor}
}

type invitationInteractor interface {
	Create(invitationRequest *dto.CreateInvitationRequest, user *model.User) (*model.Invitation, error)
	FindByInviter(inviterID int) ([]*model.Invitation, error)
	FindUserByID(id uint) (*model.User, error)
}

func toInvitationResponse(i *model.Invitation) *dto.InvitationResponse {
	return &dto.InvitationResponse{
		ID:        i.ID,
		Email:     i.Email,
		InviterID: i.InviterID,
		CreatedAt: i.CreatedAt,
	}
}

func (i *InvitationHandler) Create(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "error retrieving user from context")
	}

	request := &dto.CreateInvitationRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}

	_, err := govalidator.ValidateStruct(request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Please enter valid email")
	}

	invitation, err := i.Interactor.Create(request, user)
	if err != nil {
		return err
	}

	invitationResponse := toInvitationResponse(invitation)

	return c.JSON(http.StatusCreated, invitationResponse)
}

func (i *InvitationHandler) FindInvitationByInviter(c echo.Context) error {
	inviterID := c.QueryParam("inviter_id")
	id, err := strconv.Atoi(inviterID)
	if err != nil {
		return err
	}

	invitations, err := i.Interactor.FindByInviter(id)
	if err != nil {
		return err
	}

	invitationResponses := []*dto.InvitationResponse{}
	for _, i := range invitations {
		invitation := toInvitationResponse(i)

		invitationResponses = append(invitationResponses, invitation)
	}

	return c.JSON(http.StatusOK, invitationResponses)
}
