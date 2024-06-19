package repositories

import (
	"blog-system/internal/entities"
	"blog-system/pkg/database/postgres"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	db  *postgres.Conn
	log *logrus.Logger
}

func NewUsers(db *postgres.Conn, log *logrus.Logger) UserRepositories {
	return &userRepository{db: db, log: log}
}

func (u *userRepository) Store(ctx context.Context, p entities.User) (*entities.User, error) {
	q := `INSERT INTO users(email, full_name, password, address, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err := u.db.ExecContext(ctx, q, p.Email, p.FullName, p.Password, p.Address, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrUniqueViolation:
			return nil, fmt.Errorf("user already exists")
		default:
			return nil, err
		}
	}

	result := entities.User{
		Email:     p.Email,
		FullName:  p.FullName,
		Address:   p.Address,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	return &result, nil
}

func (u *userRepository) FindUser(ctx context.Context, email string) (*entities.User, error) {
	var (
		result entities.User
		err    error
	)

	q := `SELECT id, full_name, password, email, address FROM users WHERE email = $1`
	err = u.db.QueryRowxContext(ctx, q, email).StructScan(&result)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, fmt.Errorf("user not found")
		default:
			return nil, err
		}
	}

	return &result, nil
}

func (u *userRepository) FindUserId(ctx context.Context, id int64) (*entities.User, error) {
	var (
		result entities.User
		err    error
	)

	q := `SELECT id, full_name, email, address FROM users WHERE id = $1`
	err = u.db.QueryRowxContext(ctx, q, id).StructScan(&result)
	if err != nil {
		err = u.db.ParseSQLError(err)
		switch err {
		case postgres.ErrNoRowsFound:
			return nil, fmt.Errorf("user not found")
		default:
			return nil, err
		}
	}

	return &result, nil
}
