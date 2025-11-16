package repository

import (
	"database/sql"
	"errors"
	"log"
	"pr_reviewer/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
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
			_ = tx.Rollback()
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
        WHERE team_name = (SELECT team_name FROM users WHERE id = $1) AND id != $1 AND is_active = TRUE
        ORDER BY RANDOM() LIMIT 2
    `
	var reviewersIDs []string
	err = tx.Select(&reviewersIDs, query, authorID)

	if err != nil {
		return model.PullRequest{}, nil, err
	}

	query = `INSERT INTO pull_request_reviewers(pull_request_id, reviewer_id) VALUES ($1, $2)`

	for _, rID := range reviewersIDs {
		_, err := tx.Exec(query, id, rID)
		if err != nil {
			log.Println(err)
			return model.PullRequest{}, nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return model.PullRequest{}, nil, err
	}

	return pr, reviewersIDs, nil
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

func (r *PullRequestRepository) ReassignPullRequestReviewer(
	pullRequestID string,
	oldReviewerID string,
) (model.PullRequest, []string, string, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return model.PullRequest{}, nil, "", err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `SELECT * FROM pull_requests WHERE id=$1`

	var pr model.PullRequest
	err = tx.Get(&pr, query, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, nil, "", ErrPRNotFound
		}
		return model.PullRequest{}, nil, "", err
	}

	if pr.Status == model.PullRequestMerged {
		return model.PullRequest{}, nil, "", ErrPRMerged
	}

	query = `
		SELECT EXISTS (
			SELECT 1 FROM pull_request_reviewers
			WHERE pull_request_id=$1 AND reviewer_id=$2
		)
	`
	var exists bool
	err = tx.Get(&exists, query, pullRequestID, oldReviewerID)
	if err != nil {
		return model.PullRequest{}, nil, "", err
	}

	if !exists {
		return model.PullRequest{}, nil, "", ErrNotAssigned
	}

	var newReviewerID string
	query = `
		SELECT id
		FROM users
		WHERE team_name = (SELECT team_name FROM users WHERE id=$1)
			AND is_active = TRUE
			AND id NOT IN (
				SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id=$2
			)
		ORDER BY RANDOM()
		LIMIT 1
	`
	err = tx.Get(&newReviewerID, query, pr.AuthorID, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.PullRequest{}, nil, "", ErrNoCanditate
		}
		return model.PullRequest{}, nil, "", err
	}

	query = `
		UPDATE pull_request_reviewers
		SET reviewer_id=$1
		WHERE pull_request_id=$2 AND reviewer_id=$3
	`

	_, err = tx.Exec(query, newReviewerID, pullRequestID, oldReviewerID)
	if err != nil {
		return model.PullRequest{}, nil, "", err
	}

	query = `SELECT reviewer_id FROM pull_request_reviewers WHERE pull_request_id=$1`
	var reviewers []string
	err = tx.Select(&reviewers, query, pullRequestID)
	if err != nil {
		return model.PullRequest{}, nil, "", err
	}

	err = tx.Commit()
	if err != nil {
		return model.PullRequest{}, nil, "", err
	}

	return pr, reviewers, newReviewerID, nil
}
