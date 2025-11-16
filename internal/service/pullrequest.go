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
	pr, reviewersIDs, err := s.repo.CreatePullRequest(pullRequestID, pullRequestName, authorID)
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
		PullRequestID:     pullRequestID,
		PullRequestName:   pullRequestName,
		AuthorID:          authorID,
		Status:            "OPEN",
		AssignedReviewers: []string{},
		CreatedAt:         pr.CreatedAt,
	}
	for _, r := range reviewersIDs {
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

func (s *PullRequestService) Reassign(pullRequestID string, oldUserID string) (
	dto.ReassignPullRequestResponse, error,
) {
	pr, reviewers, replacedBy, err := s.repo.ReassignPullRequestReviewer(pullRequestID, oldUserID)
	if err != nil {
		if errors.Is(err, repository.ErrPRMerged) {
			return dto.ReassignPullRequestResponse{}, ErrPRMerged
		}

		if errors.Is(err, repository.ErrPRNotFound) {
			return dto.ReassignPullRequestResponse{}, ErrResourceNotFound
		}

		if errors.Is(err, repository.ErrNotAssigned) {
			return dto.ReassignPullRequestResponse{}, ErrNotAssigned
		}

		if errors.Is(err, repository.ErrNoCanditate) {
			return dto.ReassignPullRequestResponse{}, ErrNoCanditate
		}

		return dto.ReassignPullRequestResponse{}, err
	}

	reassignPRResponse := dto.ReassignPullRequestResponse{
		PR: dto.PullRequestShortWithReviewers{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		},
		ReplacedBy: replacedBy,
	}

	for _, r := range reviewers {
		reassignPRResponse.PR.AssignedReviewers = append(reassignPRResponse.PR.AssignedReviewers, r)
	}

	return reassignPRResponse, nil
}
