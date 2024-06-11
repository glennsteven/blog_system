package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type tagRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewTags(db *sql.DB, log *logrus.Logger) TagRepository {
	return &tagRepository{db: db, log: log}
}

func (t *tagRepository) Store(ctx context.Context, p entities.Tag) (*entities.Tag, error) {
	// Begin transaction
	tx, err := t.db.BeginTx(ctx, nil)
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

	q := `INSERT INTO tags(label) VALUES ($1) RETURNING id`

	err = t.db.QueryRowContext(ctx, q, p.Label).Scan(&p.Id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique violation
				t.log.Errorf("label must be unique: %v", err)
				return nil, fmt.Errorf("23505")
			}
		}
		t.log.Errorf("got error executing query tag: %v", err)
		return nil, err
	}

	return &p, nil
}
