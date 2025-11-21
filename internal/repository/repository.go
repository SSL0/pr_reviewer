package repository

import (
	"pr_reviewer/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	SetIsActive(userID string, isActive bool) (model.User, error)
	GetAssignedPullRequests(userID string) ([]model.PullRequest, error)
}

type TeamRepository interface {
	AddTeam(teamName string, members *[]model.User) error
	GetTeamAndMembers(teamName string) (model.Team, *[]model.User, error)
}

type PullRequestRepository interface {
	CreatePullRequest(id, name, authorID string) (model.PullRequest, []string, error)
	SetPullRequestStatus(id string, status model.PullRequestStatus) (model.PullRequest, []string, error)
	ReassignPullRequestReviewer(pullRequestID string, oldReviewerID string) (model.PullRequest, []string, string, error)
}

type Repository struct {
	UserRepository
	TeamRepository
	PullRequestRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:        NewUserRepository(db),
		TeamRepository:        NewTeamRepository(db),
		PullRequestRepository: NewPullRequestRepository(db),
	}
}
