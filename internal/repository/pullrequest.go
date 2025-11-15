package repository

import (
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type PullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) *PullRequestRepository {
	return &PullRequestRepository{db: db}
}

func (r *PullRequestRepository) CreatePullRequest(id, name, authorID string) (model.PullRequest, []string, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.Get(&exists, query, authorID)
	if err != nil {
		return model.PullRequest{}, nil, err
	}

	if !exists {
		return model.PullRequest{}, nil, ErrUserNotFound
	}

	err = r.db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id = $1)`, id)
	if err != nil {
		return model.PullRequest{}, nil, err
	}
	if exists {
		return model.PullRequest{}, nil, ErrPRExists
	}

	query = `
		INSERT INTO pull_requests(id, name, author_id, status)
		VALUES ($1, $2, $3, 'OPEN')
		RETURNING id, name, author_id, status, created_at
	`

	var pr model.PullRequest
	err = r.db.Get(&pr, query, id, name, authorID)
	if err != nil {
		return model.PullRequest{}, nil, err
	}

	query = `
		SELECT id FROM users
        WHERE team_name = (SELECT team_name FROM users WHERE id = $1) AND id != $1
        ORDER BY RANDOM() LIMIT 2
    `

	rows, err := r.db.Queryx(query, authorID)
	if err != nil {
		return model.PullRequest{}, nil, err
	}
	defer rows.Close()

	query = `INSERT INTO pull_request_reviewers(pull_request_id, reviewer_id) VALUES ($1, $2)`
	var reviewers []string
	for rows.Next() {
		var reviewerID string

		rows.Scan(&reviewerID)
		_, err := r.db.Exec(query, id, reviewerID)
		if err != nil {
			return model.PullRequest{}, nil, err
		}

		reviewers = append(reviewers, reviewerID)
	}

	return pr, reviewers, nil
}
