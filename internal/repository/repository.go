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
	GetTeamAndMembers(teamName string) (model.Team, *[]model.User, error)
}

type PullRequest interface {
	CreatePullRequest(id, name, authorID string) (model.PullRequest, []string, error)
	SetPullRequestStatus(id string, status model.PullRequestStatus) (model.PullRequest, []string, error)
	ReassignPullRequestReviewer(pullRequestID string, oldReviewerID string) (model.PullRequest, []string, string, error)
}

type Repository struct {
	User
	Team
	PullRequest
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:        NewUserRepository(db),
		Team:        NewTeamRepository(db),
		PullRequest: NewPullRequestRepository(db),
	}
}
