package mock

import (
	"github.com/stretchr/testify/mock"
	"github.com/zzokki81/uas/interactor"
	"github.com/zzokki81/uas/model"
)

type Store struct {
	mock.Mock
}

func (s *Store) CreateInvitation(invitation *model.Invitation) error {
	args := s.Called(invitation)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) CreateNotification(notification *model.Notification) error {
	args := s.Called(notification)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) Transaction() (interactor.Store, error) {
	args := s.Called()
	if args.Get(0) != nil {
		return args.Get(0).(interactor.Store), nil
	}
	return nil, args.Error(1)
}

func (s *Store) Commit() error {
	args := s.Called()
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) Rollback() {
	s.Called()
}

func (s *Store) CheckStoreConnection() error {
	args := s.Called()
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) CreateUser(user *model.User) error {
	args := s.Called(user)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) FindAllInvitations() ([]*model.Invitation, error) {
	args := s.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Invitation), nil
	}
	return nil, args.Error(1)
}

func (s *Store) FindInvitationByInviter(inviterID int) ([]*model.Invitation, error) {
	args := s.Called(inviterID)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Invitation), nil
	}
	return nil, args.Error(1)
}

func (s *Store) FindUserByID(ID uint) (*model.User, error) {
	args := s.Called(ID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), nil
	}
	return nil, args.Error(1)
}
