package interactor

import "github.com/zzokki81/uas/model"

type Store interface {
	CreateInvitation(*model.Invitation) error
	FindAllInvitations() ([]*model.Invitation, error)
	FindInvitationByInviter(int) ([]*model.Invitation, error)
	FindUserByID(uint) (*model.User, error)
	CheckStoreConnection() error
	CreateUser(*model.User) error
	Transaction() (Store, error)
	Commit() error
	Rollback()
	CreateNotification(*model.Notification) error
}
