package repository

import (
	"database/sql"
	"errors"
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SetIsActive(userID string, isActive bool) (model.User, error) {
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

func (r *UserRepository) GetAssignedPullRequests(userID string) ([]model.PullRequest, error) {
	query := `
		SELECT pr.*
		FROM pull_requests pr
		JOIN pull_request_reviewers prr
		  ON pr.id = prr.pull_request_id
		WHERE prr.reviewer_id = $1;
	`

	rows, err := r.db.Queryx(query, userID)
	if err != nil {
		return []model.PullRequest{}, err
	}
	defer rows.Close()

	result := []model.PullRequest{}

	for rows.Next() {
		var pr model.PullRequest
		if err := rows.StructScan(&pr); err != nil {
			return []model.PullRequest{}, err
		}

		result = append(result, pr)
	}

	return result, nil
}
