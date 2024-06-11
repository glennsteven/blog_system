package entities

import "time"

type Post struct {
	Id         int64     `json:"id" json:"id"`
	Title      string    `json:"title" json:"title"`
	Content    string    `json:"content" json:"content"`
	Status     int       `json:"status" json:"status"`
	Drafting   int64     `json:"drafting" json:"drafting"`
	Publishing *int64    `json:"publishing" json:"publishing,omitempty"`
	CreatedAt  time.Time `json:"created_at" json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" json:"updated_at"`
}
