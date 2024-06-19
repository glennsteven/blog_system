package repositories

import (
	"blog-system/internal/entities"
	"blog-system/pkg/database/postgres"
	"context"
	"github.com/sirupsen/logrus"
)

type userRoleRepository struct {
	db  *postgres.Conn
	log *logrus.Logger
}

func NewUserRole(db *postgres.Conn, log *logrus.Logger) RoleUser {
	return &userRoleRepository{db: db, log: log}
}

func (u *userRoleRepository) Store(ctx context.Context, p entities.UserRole) (*entities.UserRole, error) {
	q := `INSERT INTO user_roles(user_id, role_id, created_at, updated_at) VALUES ($1,$2,$3,$4) `

	_, err := u.db.ExecContext(ctx, q, p.UserId, p.RoleId, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrUniqueViolation:
			return nil, entities.ErrUserRoleAlreadyExist
		default:
			return nil, err
		}
	}

	result := entities.UserRole{
		UserId:    p.UserId,
		RoleId:    p.RoleId,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return &result, nil
}

func (u *userRoleRepository) FindUserRole(ctx context.Context, payload entities.UserRole) (*entities.UserRole, error) {
	var (
		result entities.UserRole
		err    error
	)

	q := `SELECT user_id, role_id FROM user_roles WHERE user_id = $1 AND role_id = $2`
	err = u.db.QueryRowxContext(ctx, q, payload.UserId, payload.RoleId).StructScan(&result)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, entities.ErrUserRoleNotFound
		default:
			return nil, err
		}
	}

	return &result, nil
}

func (u *userRoleRepository) FindUserIdRole(ctx context.Context, id int64) (*entities.UserRole, error) {
	var (
		result entities.UserRole
	)

	q := `SELECT user_id, role_id FROM user_roles WHERE user_id = $1 LIMIT 1`

	err := u.db.QueryRowxContext(ctx, q, id).StructScan(&result)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, entities.ErrUserRoleNotFound
		default:
			return nil, err
		}
	}

	return &result, nil
}
