package service

import (
	"errors"
	"pr_reviewer/internal/domain"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type teamService struct {
	repo *repository.Repository
}

func NewTeamService(repo *repository.Repository) *teamService {
	return &teamService{
		repo: repo,
	}
}

func (s *teamService) AddTeam(team domain.Team) (dto.Team, error) {
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

	var teamMembersDTO []dto.TeamMember

	for _, m := range team.Members {
		teamMembersDTO = append(teamMembersDTO,
			dto.TeamMember{
				UserID:   m.UserID,
				Username: m.Username,
				IsActive: m.IsActive,
			})
	}

	return dto.Team{
		TeamName: team.TeamName,
		Members:  teamMembersDTO,
	}, err
}

func (s *teamService) GetTeam(teamName string) (dto.Team, error) {
	team, users, err := s.repo.GetTeamAndMembers(teamName)

	if err != nil {
		return dto.Team{}, err
	}

	teamResponse := dto.Team{
		TeamName: team.Name,
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
