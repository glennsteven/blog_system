package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
)

type userRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewUsers(db *sql.DB, log *logrus.Logger) UserRepositories {
	return &userRepository{db: db, log: log}
}

func (u *userRepository) Store(ctx context.Context, p entities.User) (*entities.User, error) {
	var (
		result entities.User
		err    error
	)

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			return
		}
		// Commit the transaction if no error occurred
		err = tx.Commit()
	}()

	q := `INSERT INTO users(email, full_name, password, address, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *`

	_, err = u.db.ExecContext(ctx, q, p.Email, p.FullName, p.Password, p.Address, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		u.log.Errorf("got error executing query users: %v", err)
		return nil, err
	}

	err = u.db.QueryRowContext(ctx, q, p.Email, p.FullName, p.Password, p.Address, p.CreatedAt, p.UpdatedAt).Scan(
		&result.Id,
		&result.Email,
		&result.FullName,
		&result.Password,
		&result.Address,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		u.log.Errorf("got error scanning row: %v", err)
		return nil, err
	}

	return &result, nil
}

func (u *userRepository) FindUser(ctx context.Context, email string) (*entities.User, error) {
	var (
		result entities.User
		err    error
	)

	q := `SELECT full_name, email, address FROM users WHERE email = $1`
	rows, err := u.db.QueryContext(ctx, q, email)
	if err != nil {
		log.Printf("got error when find username %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.FullName, &result.Email, &result.Address)
		if err != nil {
			log.Printf("got error scan value %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
