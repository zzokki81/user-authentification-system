package model

import "time"

type Invitation struct {
	ID        uint
	Email     string
	InviterID uint
	CreatedAt time.Time
}
