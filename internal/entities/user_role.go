package entities

import "time"

type ErrorUserRole string

func (e ErrorUserRole) Error() string {
	return string(e)
}

const (
	ErrUserRoleQuery        = ErrorUserRole("Error Query")
	ErrUserRoleNotFound     = ErrorUserRole("User role not found")
	ErrUserRoleAlreadyExist = ErrorUserRole("User role already exist")
)

type UserRole struct {
	Id        int64     `db:"id" json:"id,omitempty"`
	UserId    int64     `db:"user_id" json:"user_id"`
	RoleId    int64     `db:"role_id" json:"role_id"`
	CreatedAt time.Time `cb:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
