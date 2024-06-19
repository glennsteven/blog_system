package entities

import "time"

type ErrorRole string

func (e ErrorRole) Error() string {
	return string(e)
}

const (
	ErrRoleQuery        = ErrorRole("Error Query")
	ErrRoleNotFound     = ErrorRole("Role not found")
	ErrRoleAlreadyExist = ErrorRole("Role already exist")
)

type Role struct {
	Id        int64     `db:"id" json:"id,omitempty"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
