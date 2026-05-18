package database

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
