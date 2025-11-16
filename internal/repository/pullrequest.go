package repository

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type PullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) CreatePullRequest(id, name, authorID string) (model.PullRequest, []string, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return model.PullRequest{}, nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
        INSERT INTO pull_requests(id, name, author_id, status)
        VALUES ($1, $2, $3, 'OPEN')
        RETURNING id, name, author_id, status, created_at
    `
	var pr model.PullRequest
	err = tx.Get(&pr, query, id, name, authorID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case UniqueViolationCode:
				return model.PullRequest{}, nil, ErrPRExists
			case ForeignKeyViolationCode:
				return model.PullRequest{}, nil, ErrUserNotFound
			}
		}
		return model.PullRequest{}, nil, err
	}

	query = `
        SELECT id FROM users
        WHERE team_name = (SELECT team_name FROM users WHERE id = $1) AND id != $1
        ORDER BY RANDOM() LIMIT 2
    `
	rows, err := tx.Queryx(query, authorID)
	if err != nil {
		return model.PullRequest{}, nil, err
	}
	defer rows.Close()

	insertQuery := `INSERT INTO pull_request_reviewers(pull_request_id, reviewer_id) VALUES ($1, $2)`
	var reviewers []string
	for rows.Next() {
		var reviewerID string
		if err := rows.Scan(&reviewerID); err != nil {
			return model.PullRequest{}, nil, err
		}

		_, err := tx.Exec(insertQuery, id, reviewerID)
		if err != nil {
			return model.PullRequest{}, nil, err
		}

		reviewers = append(reviewers, reviewerID)
	}

	if err := rows.Err(); err != nil {
		return model.PullRequest{}, nil, err
	}

	err = tx.Commit()
	if err != nil {
		return model.PullRequest{}, nil, err
	}

	return pr, reviewers, nil
}

func (r *PullRequestRepository) SetPullRequestStatus(id string, status model.PullRequestStatus) (model.PullRequest, []string, error) {
	query := `
		UPDATE pull_requests
		SET status = $1, merged_at = COALESCE(merged_at, NOW())
		WHERE id = $2
		RETURNING id, name, author_id, status, created_at, merged_at
	`

	var pr model.PullRequest

	err := r.db.Get(&pr, query, status, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, nil, ErrPRNotFound
		}

		return model.PullRequest{}, nil, err
	}

	query = `
		SELECT reviewer_id FROM pull_request_reviewers
	 	WHERE pull_request_id = $1
	`
	rows, err := r.db.Queryx(query, id)

	if err != nil {
		return model.PullRequest{}, nil, err
	}

	var reviewers []string

	for rows.Next() {
		var r string
		if err := rows.Scan(&r); err != nil {
			return model.PullRequest{}, nil, err
		}
		reviewers = append(reviewers, r)
	}

	return pr, reviewers, nil
}
