package entities

import "time"

type UserClaim struct {
	Id        int       `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	RoleId    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
