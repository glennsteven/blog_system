package entities

type Tag struct {
	Id    int64  `db:"id" json:"id,omitempty"`
	Label string `db:"label" json:"label"`
}
