package repositories

import (
	"blog-system/internal/entities"
	"context"
	"database/sql"
	"github.com/sirupsen/logrus"
	"time"
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

func (po *postRepository) Update(ctx context.Context, payload entities.Post, id int64) (*entities.Post, error) {
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

	q := `UPDATE posts
			SET title = $1, 
			    content = $2,
			    updated_at = $3
			WHERE id = $4 RETURNING id, title, content, updated_at`

	updatedAt := time.Now()
	row := tx.QueryRowContext(ctx, q, payload.Title, payload.Content, updatedAt, id)

	var updatedPost entities.Post
	err = row.Scan(&updatedPost.Id, &updatedPost.Title, &updatedPost.Content, &updatedPost.UpdatedAt)
	if err != nil {
		po.log.Printf("got error executing query posts: %v", err)
		return nil, err
	}

	return &updatedPost, nil
}

func (po *postRepository) FindId(ctx context.Context, id int64) (*entities.Post, error) {
	var (
		result entities.Post
		err    error
	)

	q := `SELECT id, title, content FROM posts WHERE id = $1`
	rows, err := po.db.QueryContext(ctx, q, id)
	if err != nil {
		po.log.Printf("got error when find post id %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.Title, &result.Content)
		if err != nil {
			po.log.Printf("got error scan value post %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func (po *postRepository) FindPostId(ctx context.Context, id int64) (*entities.Posts, error) {
	var (
		result entities.Posts
		err    error
	)

	q := `SELECT 
    			po.id, 
    			title, 
    			content, 
    			status, 
    			ta.label, 
    			u.full_name, 
    			u.email 
				FROM posts po
				JOIN post_tags AS pt ON pt.post_id = po.id
				JOIN tags AS ta ON ta.id = pt.tag_id
				JOIN users AS u ON u.id = po.drafting
				WHERE po.id = $1 `

	rows, err := po.db.QueryContext(ctx, q, id)
	if err != nil {
		po.log.Printf("got error when find post id %v", err)
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&result.Id, &result.Title, &result.Content, &result.Status, &result.Label, &result.FullName, &result.Email)
		if err != nil {
			po.log.Printf("got error scan value post %v", err)
			return nil, err
		}
		return &result, nil
	} else {
		return nil, nil
	}
}

func (po *postRepository) DeletePost(ctx context.Context, id int64) error {
	// Begin transaction
	tx, err := po.db.BeginTx(ctx, nil)
	if err != nil {
		po.log.Printf("could not begin transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	q := `DELETE FROM posts WHERE id = $1`

	// Execute the delete query
	_, err = tx.ExecContext(ctx, q, id)
	if err != nil {
		po.log.Printf("process destroy post got error: %v", err)
		return err
	}

	return nil
}
