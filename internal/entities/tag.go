package entities

type Tag struct {
	Id    int64  `db:"id" json:"id"`
	Label string `db:"label" json:"label"`
}
