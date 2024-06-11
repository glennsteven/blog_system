package entities

import "time"

type UserClaim struct {
	Id        int       `db:"id" json:"id,omitempty"`
	Email     string    `db:"email" json:"email,omitempty"`
	FullName  string    `db:"full_name" json:"full_name,omitempty"`
	Password  string    `db:"password" json:"password,omitempty"`
	Address   string    `db:"address" json:"address,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
