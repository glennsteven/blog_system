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

	q := `INSERT INTO roles(name, created_at, updated_at) VALUES ($1,$2,$3)`

	_, err = r.db.ExecContext(ctx, q, p.Name, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		r.log.Errorf("got error executing query roles: %v", err)
		return nil, err
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
	rows, err := r.db.QueryContext(ctx, q, name)
	if err != nil {
		log.Printf("got error when find role name %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Name)
		if err != nil {
			log.Printf("got error scan value role %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func (r *roleRepository) FindRoleId(ctx context.Context, id int64) (*entities.Role, error) {
	var (
		result entities.Role
		err    error
	)

	q := `SELECT id, name FROM roles WHERE id = $1`
	rows, err := r.db.QueryContext(ctx, q, id)
	if err != nil {
		log.Printf("got error when find role id %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.Name)
		if err != nil {
			log.Printf("got error scan value role %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
