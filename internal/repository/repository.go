package repository

import (
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type User interface {
	SetIsActive(userID string, isActive bool) (model.User, error)
	GetAssignedPullRequests(userID string) ([]model.PullRequest, error)
}

type Team interface {
	AddTeam(teamName string, members *[]model.User) error
	GetTeamMembers(teamName string) (*[]model.User, error)
}

type Repository struct {
	User
	Team
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
		Team: NewTeamRepository(db),
	}
}
