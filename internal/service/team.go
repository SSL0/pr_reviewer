package service

import (
	"pr_service/internal/model"
	"pr_service/internal/repository"
)

type TeamService struct {
	repo *repository.Repository
}

func NewTeamService(repo *repository.Repository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}

func (s *TeamService) Add(team model.Team) (model.Team, error) {
	return model.Team{}, nil
}

func (s *TeamService) Get(teamName string) (model.Team, error) {
	return model.Team{}, nil
}
