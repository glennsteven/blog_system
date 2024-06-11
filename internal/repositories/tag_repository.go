package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
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
				t.log.Infof("label must be unique: %v", err)
				return nil, fmt.Errorf("23505")
			}
		}
		t.log.Errorf("got error executing query tag: %v", err)
		return nil, err
	}

	return &p, nil
}

func (t *tagRepository) Update(ctx context.Context, payload entities.Tag, label string) (*entities.Tag, error) {
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

	q := `UPDATE tags
			SET label = $1
			WHERE label = $4 RETURNING id, lable`

	row := tx.QueryRowContext(ctx, q, payload.Label, label)

	var updatedTag entities.Tag
	err = row.Scan(&updatedTag.Id, &updatedTag.Label)
	if err != nil {
		log.Printf("got error executing query tags: %v", err)
		return nil, err
	}

	return &updatedTag, nil
}

func (t *tagRepository) FindLabel(ctx context.Context, label string) (*entities.Tag, error) {
	var (
		result entities.Tag
		err    error
	)

	q := `SELECT id, label FROM tags WHERE label = $1`
	rows, err := t.db.QueryContext(ctx, q, label)
	if err != nil {
		log.Printf("got error when find tag label %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.Label)
		if err != nil {
			log.Printf("got error scan value tag %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}
