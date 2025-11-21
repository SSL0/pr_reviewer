package service

import (
	"pr_reviewer/internal/domain"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type TeamService interface {
	AddTeam(team domain.Team) (dto.Team, error)
	GetTeam(teamName string) (dto.Team, error)
}

type UserService interface {
	SetUserIsActive(userID string, isActive bool) (model.User, error)
	GetUserReviews(userID string) (dto.UserReviewResponse, error)
}

type PullRequestService interface {
	CreatePullRequest(pullRequestID, pullRequestName, authorID string) (dto.PullRequest, error)
	MergePullRequest(pullReqeustID string) (dto.MergePullRequestResponse, error)
	ReassignPullRequestReviewer(pullRequestID string, oldUserID string) (dto.ReassignPullRequestResponse, error)
}

type Service struct {
	TeamService
	UserService
	PullRequestService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		TeamService:        NewTeamService(repo),
		UserService:        NewUserService(repo),
		PullRequestService: NewPullRequestService(repo),
	}
}
