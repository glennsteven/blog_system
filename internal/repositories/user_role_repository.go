package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
)

type userRoleRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewUserRole(db *sql.DB, log *logrus.Logger) RoleUser {
	return &userRoleRepository{db: db, log: log}
}

func (u *userRoleRepository) Store(ctx context.Context, p entities.UserRole) (*entities.UserRole, error) {

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

	q := `INSERT INTO user_roles(user_id, role_id, created_at, updated_at) VALUES ($1,$2,$3,$4) `

	_, err = u.db.ExecContext(ctx, q, p.UserId, p.RoleId, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		u.log.Errorf("got error executing query role user: %v", err)
		return nil, err
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
	rows, err := u.db.QueryContext(ctx, q, payload.UserId, payload.RoleId)
	if err != nil {
		log.Printf("got error when find role user %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.UserId, &result.RoleId)
		if err != nil {
			log.Printf("got error scan value %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func (u *userRoleRepository) FindUserIdRole(ctx context.Context, id int64) (*entities.UserRole, error) {
	var (
		result entities.UserRole
	)

	q := `SELECT user_id, role_id FROM user_roles WHERE user_id = $1 LIMIT 1`

	rows, err := u.db.QueryContext(ctx, q, id)
	if err != nil {
		log.Printf("got error when find role user %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.UserId, &result.RoleId)
		if err != nil {
			log.Printf("got error scan value %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
