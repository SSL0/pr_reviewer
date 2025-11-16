package service

import (
	"errors"
	"pr_reviewer/internal/domain"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

type TeamService struct {
	repo *repository.Repository
}

func NewTeamService(repo *repository.Repository) *TeamService {
	return &TeamService{
		repo: repo,
	}
}

func (s *TeamService) AddTeam(team domain.Team) (dto.Team, error) {
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

	var TeamMembersDTO []dto.TeamMember

	for _, m := range team.Members {
		TeamMembersDTO = append(TeamMembersDTO,
			dto.TeamMember{
				UserID:   m.UserID,
				Username: m.Username,
				IsActive: m.IsActive,
			})
	}

	return dto.Team{
		TeamName: team.TeamName,
		Members:  TeamMembersDTO,
	}, err
}

func (s *TeamService) GetTeam(teamName string) (dto.Team, error) {
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
