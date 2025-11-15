package service

import (
	"pr_service/internal/model"
	"pr_service/internal/repository"
)

type Team interface {
	Add(team model.Team) (model.Team, error)
	Get(teamName string) (model.Team, error)
}

type User interface {
	SetIsActive(userID string, isActive bool) (model.User, error)
	GetReview(userID string) (string, model.PullRequestShort)
}

type PullReqeust interface {
	Create(pullRequestID string, pullRequestName string, authorID string) (model.PullRequest, error)
	Merge(pullReqeustID string) (model.PullRequest, error)
	Reassign(pullRequestID string, oldUserID string) (model.PullRequest, string, error)
}

type Service struct {
	Team
	User
	PullReqeust
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Team:        NewTeamService(repo),
		User:        NewUserService(repo),
		PullReqeust: NewPullRequestService(repo),
	}
}
