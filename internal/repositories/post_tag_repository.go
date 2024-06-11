package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type postTagRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewPostTag(db *sql.DB, log *logrus.Logger) PostTagRepository {
	return &postTagRepository{db: db, log: log}
}

func (ps *postTagRepository) Store(ctx context.Context, p entities.PostTag) error {
	// Begin transaction
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
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

	q := `INSERT INTO post_tags(post_id, tag_id) VALUES ($1,$2)`

	_, err = ps.db.ExecContext(ctx, q, p.PostId, p.TagId)
	if err != nil {
		ps.log.Errorf("got error executing query post tag: %v", err)
		return err
	}

	return nil
}
