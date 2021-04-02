package interactor

import (
	"github.com/zzokki81/uas/model"
)

type userStore interface {
	FindUserByID(uint) (*model.User, error)
}

type User struct {
	store userStore
}

func NewUser(store userStore) *User {
	return &User{
		store: store,
	}
}

func (u *User) FindUserByID(id uint) (*model.User, error) {
	return u.store.FindUserByID(id)
}
