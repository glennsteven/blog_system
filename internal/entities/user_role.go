package entities

type UserRole struct {
	Id     int64 `db:"id" json:"id"`
	UserId int64 `db:"user_id" json:"user_id"`
	RoleId int64 `db:"role_id" json:"role_id"`
}
