package service

import (
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
	return model.User{}, nil
}

func (s *UserService) GetReview(userID string) (string, model.PullRequestShort) {
	return "", model.PullRequestShort{}
}
