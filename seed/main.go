package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zzokki81/uas/model"
	"github.com/zzokki81/uas/store/postgres"
)

var users = []*model.User{
	{
		ID:        1,
		Email:     "john@kalhun.com",
		Name:      "John",
		CreatedAt: time.Time{},
	},
}

func areTablesEmpty(db *gorm.DB) error {
	tables := []string{"users"}
	for _, table := range tables {
		var count int

		if err := db.Table(table).Count(&count).Error; err != nil {
			return err
		}
		if count != 0 {
			return errors.New("table " + table + " is not empty")
		}
	}
	return nil
}

func main() {

	config := postgres.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "uas",
	}

	store, err := postgres.Open(config)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	if err = areTablesEmpty(store.DB()); err != nil {
		panic(err)
	}

	tx, err := store.Transaction()
	if err != nil {
		panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error))
			tx.Rollback()
		}
	}()

	for _, user := range users {
		if err := tx.CreateUser(user); err != nil {
			panic(err)
		}
	}

	tx.Commit()
}
