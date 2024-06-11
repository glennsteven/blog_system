package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"log"
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

func (ps *postTagRepository) FindPostId(ctx context.Context, postId int64) ([]entities.PostTag, error) {
	var result []entities.PostTag

	q := `SELECT post_id, tag_id FROM post_tags WHERE post_id = $1`
	rows, err := ps.db.QueryContext(ctx, q, postId)
	if err != nil {
		log.Printf("got error when finding post id %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postTag entities.PostTag
		err = rows.Scan(&postTag.PostId, &postTag.TagId)
		if err != nil {
			log.Printf("got error scanning value post tag %v", err)
			return nil, err
		}
		result = append(result, postTag)
	}

	// Check for any error that may have occurred during iteration
	if err = rows.Err(); err != nil {
		log.Printf("got error iterating rows %v", err)
		return nil, err
	}

	return result, nil
}
