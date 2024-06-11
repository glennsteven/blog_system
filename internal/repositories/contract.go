package repositories

import (
	"blog-system/internal/entities"
	"context"
)

type UserRepositories interface {
	Store(ctx context.Context, payload entities.User) (*entities.User, error)
	FindUser(ctx context.Context, email string) (*entities.User, error)
}

type RoleRepositories interface {
	Store(ctx context.Context, payload entities.Role) (*entities.Role, error)
	FindRole(ctx context.Context, name string) (*entities.Role, error)
}
