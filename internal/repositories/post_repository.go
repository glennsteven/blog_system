package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
)

type postRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewPost(db *sql.DB, log *logrus.Logger) PostRepository {
	return &postRepository{db: db, log: log}
}

func (po *postRepository) Store(ctx context.Context, p entities.Post) (*entities.Post, error) {
	// Begin transaction
	tx, err := po.db.BeginTx(ctx, nil)
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

	q := `INSERT INTO posts(title, content, status, drafting, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`

	err = po.db.QueryRowContext(ctx, q, p.Title, p.Content, p.Status, p.Drafting, p.CreatedAt, p.UpdatedAt).Scan(&p.Id)
	if err != nil {
		po.log.Errorf("got error executing query post: %v", err)
		return nil, err
	}

	return &p, nil
}
