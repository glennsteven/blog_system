package repositories

import (
	"blog-system/internal/entities"
	"context"
)

type UserRepositories interface {
	Store(ctx context.Context, payload entities.User) (*entities.User, error)
	FindUser(ctx context.Context, email string) (*entities.User, error)
	FindUserId(ctx context.Context, id int64) (*entities.User, error)
}

type RoleRepositories interface {
	Store(ctx context.Context, payload entities.Role) (*entities.Role, error)
	FindRole(ctx context.Context, name string) (*entities.Role, error)
	FindRoleId(ctx context.Context, id int64) (*entities.Role, error)
}

type RoleUser interface {
	Store(ctx context.Context, payload entities.UserRole) (*entities.UserRole, error)
	FindUserRole(ctx context.Context, payload entities.UserRole) (*entities.UserRole, error)
	FindUserIdRole(ctx context.Context, id int64) (*entities.UserRole, error)
}

type PostRepository interface {
	Store(ctx context.Context, payload entities.Post) (*entities.Post, error)
}

type TagRepository interface {
	Store(ctx context.Context, payload entities.Tag) (*entities.Tag, error)
}

type PostTagRepository interface {
	Store(ctx context.Context, payload entities.PostTag) error
}
