package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

func (s *Store) DB() *gorm.DB {
	return s.db
}
func (store *Store) Commit() error {
	return store.db.Commit().Error
}

func (store *Store) Rollback() {
	store.db.Rollback()
}
