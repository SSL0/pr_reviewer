package service

import (
	"errors"
	"pr_service/internal/dto"
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

func (s *TeamService) Add(team dto.Team) (dto.Team, error) {
	var teamUsers []model.User

	for _, m := range team.Members {
		teamUsers = append(
			teamUsers,
			model.User{
				ID:       m.UserID,
				Username: m.Username,
				IsActive: m.IsActive,
			},
		)
	}

	err := s.repo.AddTeam(team.TeamName, &teamUsers)
	if errors.Is(err, repository.ErrTeamExists) {
		return dto.Team{}, ErrTeamExists
	}

	return team, err
}

func (s *TeamService) Get(teamName string) (dto.Team, error) {
	users, err := s.repo.GetTeamMembers(teamName)

	if err != nil {
		return dto.Team{}, err
	}

	if len(*users) == 0 {
		return dto.Team{}, ErrTeamNotFound
	}

	teamResponse := dto.Team{
		TeamName: teamName,
		Members:  []dto.TeamMember{},
	}

	for _, m := range *users {
		member := dto.TeamMember{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
		teamResponse.Members = append(teamResponse.Members, member)
	}

	return teamResponse, nil
}
