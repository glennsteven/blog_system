package repositories

import (
	"blog-system/internal/entities"
	"blog-system/pkg/database/postgres"
	"context"
	"github.com/sirupsen/logrus"
)

type roleRepository struct {
	db  *postgres.Conn
	log *logrus.Logger
}

func NewRoles(db *postgres.Conn, log *logrus.Logger) RoleRepositories {
	return &roleRepository{db: db, log: log}
}

func (r *roleRepository) Store(ctx context.Context, p entities.Role) (*entities.Role, error) {
	q := `INSERT INTO roles(name, created_at, updated_at) VALUES ($1,$2,$3)`

	_, err := r.db.ExecContext(ctx, q, p.Name, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		err = r.db.ParseSQLError(err)
		switch err {
		case postgres.ErrUniqueViolation:
			return nil, entities.ErrRoleAlreadyExist
		default:
			return nil, err
		}
	}

	result := entities.Role{
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return &result, nil
}

func (r *roleRepository) FindRole(ctx context.Context, name string) (*entities.Role, error) {
	var (
		result entities.Role
		err    error
	)

	q := `SELECT name FROM roles WHERE name = $1`
	err = r.db.QueryRowxContext(ctx, q, name).StructScan(&result)
	if err != nil {
		err = r.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, entities.ErrRoleNotFound
		default:
			return nil, err
		}
	}

	return &result, nil
}

func (r *roleRepository) FindRoleId(ctx context.Context, id int64) (*entities.Role, error) {
	var (
		result entities.Role
		err    error
	)

	q := `SELECT id, name FROM roles WHERE id = $1`
	err = r.db.QueryRowxContext(ctx, q, id).StructScan(&result)
	if err != nil {
		err = r.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, entities.ErrRoleNotFound
		default:
			return nil, err
		}
	}

	return &result, nil
}
