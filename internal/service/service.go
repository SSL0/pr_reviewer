package service

import (
	"pr_reviewer/internal/domain"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type Team interface {
	AddTeam(team domain.Team) (dto.Team, error)
	GetTeam(teamName string) (dto.Team, error)
}

type User interface {
	SetUserIsActive(userID string, isActive bool) (model.User, error)
	GetUserReviews(userID string) (dto.UserReviewResponse, error)
}

type PullReqeust interface {
	CreatePullRequest(pullRequestID, pullRequestName, authorID string) (dto.PullRequest, error)
	MergePullRequest(pullReqeustID string) (dto.MergePullRequestResponse, error)
	ReassignPullRequestReviewer(pullRequestID string, oldUserID string) (dto.ReassignPullRequestResponse, error)
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
