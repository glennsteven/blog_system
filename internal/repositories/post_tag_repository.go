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

	if err = rows.Err(); err != nil {
		log.Printf("got error iterating rows %v", err)
		return nil, err
	}

	return result, nil
}

func (ps *postTagRepository) FindTagId(ctx context.Context, tagId int64) ([]entities.PostTag, error) {
	var result []entities.PostTag

	q := `SELECT post_id, tag_id FROM post_tags WHERE tag_id = $1`
	rows, err := ps.db.QueryContext(ctx, q, tagId)
	if err != nil {
		log.Printf("got error when finding tag id %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postTag entities.PostTag
		err = rows.Scan(&postTag.PostId, &postTag.TagId)
		if err != nil {
			log.Printf("got error scanning value tag post %v", err)
			return nil, err
		}
		result = append(result, postTag)
	}

	if err = rows.Err(); err != nil {
		log.Printf("got error iterating rows %v", err)
		return nil, err
	}

	return result, nil
}

func (ps *postTagRepository) DeletePostTag(ctx context.Context, id int64) error {
	// Begin transaction
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		ps.log.Printf("could not begin transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	q := `DELETE FROM post_tags WHERE post_id = $1`

	_, err = tx.ExecContext(ctx, q, id)
	if err != nil {
		ps.log.Printf("process destroy post tag got error: %v", err)
		return err
	}

	return nil
}
