package model

import (
	"time"
)

const (
	AdminStatusActive  = "active"
	AdminStatusBlocked = "blocked"
)

type Admin struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Status    string    `db:"status"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt time.Time `db:"created_at"`
}
