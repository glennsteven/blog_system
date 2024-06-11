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

	q := `INSERT INTO users(email, full_name, password, address, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`

	_, err = u.db.ExecContext(ctx, q, p.Email, p.FullName, p.Password, p.Address, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		u.log.Errorf("got error executing query users: %v", err)
		return nil, err
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
	rows, err := u.db.QueryContext(ctx, q, email)
	if err != nil {
		log.Printf("got error when find username %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.FullName, &result.Password, &result.Email, &result.Address)
		if err != nil {
			log.Printf("got error scan value %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func (u *userRepository) FindUserId(ctx context.Context, id int64) (*entities.User, error) {
	var (
		result entities.User
		err    error
	)

	q := `SELECT id, full_name, email, address FROM users WHERE id = $1`
	rows, err := u.db.QueryContext(ctx, q, id)
	if err != nil {
		log.Printf("got error when find user id %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.FullName, &result.Email, &result.Address)
		if err != nil {
			log.Printf("got error scan value %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
