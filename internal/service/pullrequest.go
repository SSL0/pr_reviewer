package service

import (
	"errors"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type PullRequestService struct {
	repo *repository.Repository
}

func NewPullRequestService(repo *repository.Repository) *PullRequestService {
	return &PullRequestService{
		repo: repo,
	}
}

func (s *PullRequestService) Create(pullRequestID, pullRequestName, authorID string) (dto.PullRequest, error) {
	pr, reviewers, err := s.repo.CreatePullRequest(pullRequestID, pullRequestName, authorID)
	if err != nil {

		if errors.Is(err, repository.ErrPRExists) {
			return dto.PullRequest{}, ErrPRExists
		}

		if errors.Is(err, repository.ErrUserNotFound) {
			return dto.PullRequest{}, ErrResourceNotFound
		}

		return dto.PullRequest{}, err
	}

	dtoPR := dto.PullRequest{
		PullRequestID:   pullRequestID,
		PullRequestName: pullRequestName,
		AuthorID:        authorID,
		Status:          "OPEN",
		CreatedAt:       pr.CreatedAt,
	}
	for _, r := range reviewers {
		dtoPR.AssignedReviewers = append(dtoPR.AssignedReviewers, r)
	}
	return dtoPR, nil
}

func (s *PullRequestService) Merge(pullReqeustID string) (dto.MergePullRequestResponse, error) {
	pr, reviewers, err := s.repo.SetPullRequestStatus(pullReqeustID, model.PullRequestMerged)

	if err != nil {
		if errors.Is(err, repository.ErrPRNotFound) {
			return dto.MergePullRequestResponse{}, ErrResourceNotFound
		}
		return dto.MergePullRequestResponse{}, err
	}

	mergePRResponse := dto.MergePullRequestResponse{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
		MergedAt:        pr.MergedAt,
	}

	for _, r := range reviewers {
		mergePRResponse.AssignedReviewers = append(mergePRResponse.AssignedReviewers, r)
	}

	return mergePRResponse, nil
}

func (s *PullRequestService) Reassign(pullRequestID string, oldUserID string) (model.PullRequest, string, error) {
	return model.PullRequest{}, "", nil
}
