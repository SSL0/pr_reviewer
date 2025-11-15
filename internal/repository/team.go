package repository

import (
	"pr_reviewer/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) AddTeam(teamName string, members *[]model.User) error {
	const uniqueViolationCode = "23505"

	query := `INSERT INTO teams(name) VALUES ($1)`

	_, err := r.db.Exec(query, teamName)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == uniqueViolationCode {
			return ErrTeamExists
		}
		return err
	}

	userQuery := `
		INSERT INTO users(id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET team_name = EXCLUDED.team_name,
		    is_active = EXCLUDED.is_active,
		    username = EXCLUDED.username;
	`

	for _, m := range *members {
		_, err := r.db.Exec(userQuery, m.ID, m.Username, teamName, m.IsActive)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) GetTeamMembers(teamName string) (*[]model.User, error) {
	query := `
		SELECT * FROM users
		WHERE team_name = $1;
	`

	var result []model.User
	rows, err := r.db.Queryx(query, teamName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user model.User
		err := rows.StructScan(&user)

		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}

	return &result, nil
}
