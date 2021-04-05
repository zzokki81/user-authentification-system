package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/zzokki81/uas/handler/dto"
	"github.com/zzokki81/uas/model"
)

type Invitation struct {
	mock.Mock
}

func (i *Invitation) Create(request *dto.CreateInvitationRequest, user *model.User) (*model.Invitation, error) {
	args := i.Called(request, user)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Invitation), nil
	}
	return nil, args.Error(1)
}

func (i *Invitation) FindByInviter(inviterID int) ([]*model.Invitation, error) {
	args := i.Called(inviterID)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Invitation), nil
	}

	return nil, args.Error(1)
}

func (i *Invitation) FindUserByID(id uint) (*model.User, error) {
	args := i.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), nil
	}
	return nil, args.Error(1)
}
