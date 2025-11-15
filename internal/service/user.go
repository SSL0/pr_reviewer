package service

import (
	"database/sql"
	"pr_service/internal/dto"
	"pr_service/internal/model"
	"pr_service/internal/repository"
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
	user, err := s.repo.SetIsActiveByUserID(userID, isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, ErrUserNotFound
		}

		return model.User{}, err
	}

	return user, nil
}

func (s *UserService) GetReview(userID string) (dto.GetReviewResponse, error) {
	prs, err := s.repo.GetAssignedPullRequestsByUserID(userID)

	if err != nil {
		return dto.GetReviewResponse{}, err
	}

	result := dto.GetReviewResponse{
		UserID:       userID,
		PullRequests: []dto.PullRequestShort{},
	}

	for _, pr := range prs {
		result.PullRequests = append(
			result.PullRequests,
			dto.PullRequestShort{
				PullRequestID:   pr.ID,
				PullRequestName: pr.Name,
				AuthorID:        pr.AuthorID,
				Status:          pr.Status,
			},
		)
	}

	return result, err
}
