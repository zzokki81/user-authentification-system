package model

import "time"

type User struct {
	ID        uint
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
