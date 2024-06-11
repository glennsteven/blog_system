package requests

type AssignRoleRequest struct {
	UserId int64 `json:"user_id"`
	RoleId int64 `json:"role_id"`
}
