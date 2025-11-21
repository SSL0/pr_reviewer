package repository

import (
	"context"
	"database/sql"
	"errors"
	"pr_reviewer/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type teamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *teamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) AddTeam(teamName string, members *[]model.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `INSERT INTO teams(name) VALUES ($1)`
	_, err = tx.ExecContext(context.Background(), query, teamName)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationCode {
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
		_, err = tx.ExecContext(context.Background(), userQuery, m.ID, m.Username, teamName, m.IsActive)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *teamRepository) GetTeamAndMembers(teamName string) (model.Team, *[]model.User, error) {
	query := `SELECT name FROM teams WHERE name = $1`

	var team model.Team
	err := r.db.Get(&team, query, teamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Team{}, nil, ErrTeamNotFound
		}

		return model.Team{}, nil, err
	}

	query = `
		SELECT id, username, team_name, is_active FROM users
		WHERE team_name = $1;
	`

	var result []model.User
	err = r.db.Select(&result, query, teamName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.Team{}, nil, err
	}

	return team, &result, nil
}
