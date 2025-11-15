package service

import (
	"pr_service/internal/model"
	"pr_service/internal/repository"
)

type PullRequestService struct {
	repo *repository.Repository
}

func NewPullRequestService(repo *repository.Repository) *PullRequestService {
	return &PullRequestService{
		repo: repo,
	}
}

func (s *PullRequestService) Create(pullRequestID string, pullRequestName string, authorID string) (model.PullRequest, error) {
	return model.PullRequest{}, nil
}

func (s *PullRequestService) Merge(pullReqeustID string) (model.PullRequest, error) {
	return model.PullRequest{}, nil
}
func (s *PullRequestService) Reassign(pullRequestID string, oldUserID string) (model.PullRequest, string, error) {
	return model.PullRequest{}, "", nil
}
