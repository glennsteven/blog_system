package entities

import "time"

type Post struct {
	Id         int64
	Title      string
	Content    string
	Status     int
	Drafting   int64
	Publishing *int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
