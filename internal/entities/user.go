package entities

import "time"

type ErrorUser string

func (e ErrorUser) Error() string {
	return string(e)
}

const (
	ErrUserQuery        = ErrorUser("Error Query")
	ErrUserNotFound     = ErrorUser("User not found")
	ErrUserAlreadyExist = ErrorUser("User already exist")
)

type User struct {
	Id        int64     `db:"id" json:"id,omitempty"`
	Email     string    `db:"email" json:"email"`
	FullName  string    `db:"full_name" json:"full_name"`
	Password  string    `db:"password" json:"password,omitempty"`
	Address   string    `db:"address" json:"address"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type UserAuth struct {
	Id       int64  `db:"id" json:"id,omitempty"`
	Email    string `db:"email" json:"email"`
	FullName string `db:"full_name" json:"full_name"`
	RoleId   int64  `db:"role_id" json:"role_id"`
}
