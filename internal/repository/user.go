package repository

import (
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(userID string) (model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	var user model.User

	err := r.db.QueryRowx(query, userID).StructScan(user)
	return user, err
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

	return result, err
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

	result := []model.PullRequest{}

	for rows.Next() {
		var pr model.PullRequest
		err := rows.StructScan(&pr)

		if err != nil {
			return []model.PullRequest{}, err
		}

		result = append(result, pr)
	}

	return result, nil
}
