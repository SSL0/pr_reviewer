package repository

import (
	"pr_service/internal/model"

	"github.com/jmoiron/sqlx"
)

type User interface {
	SetIsActiveByUserID(userID string, isActive bool) (model.User, error)
	GetAssignedPullRequestsByUserID(userID string) ([]model.PullRequest, error)
}

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
