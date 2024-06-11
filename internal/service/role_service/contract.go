package role_service

import (
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
)

type RoleUser interface {
	RoleUser(ctx context.Context, payload requests.RoleRequest) (resources.Response, error)
	AssignRole(ctx context.Context, payload requests.AssignRoleRequest) (resources.Response, error)
}
