package service

import (
	"database/sql"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SetIsActive(userID string, isActive bool) (model.User, error) {
	user, err := s.repo.SetIsActive(userID, isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetReview(userID string) (dto.UserReviewResponse, error) {
	prs, err := s.repo.GetAssignedPullRequests(userID)

	if err != nil {
		return dto.UserReviewResponse{}, err
	}

	getReviewResponse := dto.UserReviewResponse{
		UserID:       userID,
		PullRequests: []dto.PullRequestShort{},
	}

	for _, pr := range prs {
		getReviewResponse.PullRequests = append(
			getReviewResponse.PullRequests,
			dto.PullRequestShort{
				PullRequestID:   pr.ID,
				PullRequestName: pr.Name,
				AuthorID:        pr.AuthorID,
				Status:          pr.Status,
			},
		)
	}

	return getReviewResponse, err
}
