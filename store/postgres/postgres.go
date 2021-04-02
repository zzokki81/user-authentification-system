package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zzokki81/uas/interactor"
	"github.com/zzokki81/uas/model"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func Open(config PostgresConfig) (*Store, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	store := NewStore(db)
	return store, nil
}

func (store *Store) CheckStoreConnection() error {
	return store.db.DB().Ping()
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) CreateUser(user *model.User) error {
	return store.db.Create(&user).Error
}

func (store *Store) FindInvitationByInviter(inviterID int) ([]*model.Invitation, error) {
	invitations := []*model.Invitation{}
	if err := store.db.Where("inviter_id = ?", inviterID).Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

func (store *Store) CreateInvitation(invitation *model.Invitation) error {
	return store.db.Create(invitation).Error

}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (store *Store) Transaction() (interactor.Store, error) {
	tx := store.db.Begin()
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return NewStore(tx), nil
}

func (store *Store) Commit() error {
	return store.db.Commit().Error
}

func (store *Store) Rollback() {
	store.db.Rollback()
}

func (s *Store) FindUserByID(id uint) (*model.User, error) {
	user := &model.User{}
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (store *Store) FindAllInvitations() ([]*model.Invitation, error) {
	invitations := []*model.Invitation{}
	if err := store.db.Find(&invitations).Error; err != nil {
		return nil, err
	}
	return invitations, nil
}

func (store *Store) CreateNotification(notification *model.Notification) error {
	return store.db.Create(notification).Error
}
