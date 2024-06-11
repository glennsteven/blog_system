package entities

import "time"

type User struct {
	Id        int64
	Email     string
	FullName  string
	Password  string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
