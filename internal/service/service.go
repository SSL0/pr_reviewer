package service

import (
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type Team interface {
	Add(team dto.Team) (dto.Team, error)
	Get(teamName string) (dto.Team, error)
}

type User interface {
	SetIsActive(userID string, isActive bool) (model.User, error)
	GetReview(userID string) (dto.UserReviewResponse, error)
}

type PullReqeust interface {
	Create(pullRequestID, pullRequestName, authorID string) (dto.PullRequest, error)
	Merge(pullReqeustID string) (dto.MergePullRequestResponse, error)
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
