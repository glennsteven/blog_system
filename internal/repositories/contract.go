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
	Update(ctx context.Context, payload entities.Post, id int64) (*entities.Post, error)
	FindId(ctx context.Context, id int64) (*entities.Post, error)
	FindPostId(ctx context.Context, id int64) (*entities.Posts, error)
	DeletePost(ctx context.Context, id int64) error
}

type TagRepository interface {
	Store(ctx context.Context, payload entities.Tag) (*entities.Tag, error)
	Update(ctx context.Context, payload entities.Tag, label string) (*entities.Tag, error)
	FindLabel(ctx context.Context, label string) (*entities.Tag, error)
	FindId(ctx context.Context, id int64) (*entities.Tag, error)
}

type PostTagRepository interface {
	Store(ctx context.Context, payload entities.PostTag) error
	FindPostId(ctx context.Context, postId int64) ([]entities.PostTag, error)
	FindTagId(ctx context.Context, tagId int64) ([]entities.PostTag, error)
	DeletePostTag(ctx context.Context, id int64) error
}
