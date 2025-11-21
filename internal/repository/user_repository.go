package repository

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) SetIsActive(userID string, isActive bool) (model.User, error) {
	query := `
		UPDATE users
		SET is_active = $1
		WHERE id = $2
		RETURNING id, username, team_name, is_active
	`

	var result model.User
	err := r.db.Get(&result, query, isActive, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, err
	}

	return result, nil
}

func (r *userRepository) GetAssignedPullRequests(userID string) ([]model.PullRequest, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exsits bool
	err := r.db.Get(&exsits, query, userID)

	if err != nil {
		return nil, err
	}

	if !exsits {
		return nil, ErrUserNotFound
	}

	query = `
		SELECT pr.*
		FROM pull_requests pr
		JOIN pull_request_reviewers prr
		  ON pr.id = prr.pull_request_id
		WHERE prr.reviewer_id = $1;
	`
	var result []model.PullRequest
	err = r.db.Select(&result, query, userID)

	if err != nil {
		return []model.PullRequest{}, err
	}

	return result, nil
}
