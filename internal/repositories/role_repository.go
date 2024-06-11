package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
)

type roleRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewRoles(db *sql.DB, log *logrus.Logger) RoleRepositories {
	return &roleRepository{db: db, log: log}
}

func (r *roleRepository) Store(ctx context.Context, p entities.Role) (*entities.Role, error) {
	var (
		result entities.Role
		err    error
	)

	// Begin transaction
	tx, err := r.db.BeginTx(ctx, nil)
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

	q := `INSERT INTO roles(name, created_at, updated_at) VALUES ($1,$2,$3) RETURNING *`

	_, err = r.db.ExecContext(ctx, q, p.Name, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		r.log.Errorf("got error executing query roles: %v", err)
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, q, p.Name, p.CreatedAt, p.UpdatedAt).Scan(
		&result.Id,
		&result.Name,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		r.log.Errorf("got error scanning row role_controller: %v", err)
		return nil, err
	}

	return &result, nil
}

func (r *roleRepository) FindRole(ctx context.Context, name string) (*entities.Role, error) {
	var (
		result entities.Role
		err    error
	)

	q := `SELECT name FROM roles WHERE name = $1`
	rows, err := r.db.QueryContext(ctx, q, name)
	if err != nil {
		log.Printf("got error when find role_controller name %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Name)
		if err != nil {
			log.Printf("got error scan value role_controller %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
