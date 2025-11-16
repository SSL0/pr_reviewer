package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestTeamRepository(t *testing.T) {
	testDB, err := repository.NewPostgres(dbDSN)
	require.NoError(t, err)

	repo := repository.NewRepository(testDB)

	teamName := "team11231324"
	members := []model.User{
		{ID: "u1", Username: "Alice", TeamName: teamName, IsActive: true},
		{ID: "u2", Username: "Bob", TeamName: teamName, IsActive: false},
	}

	t.Run("AddTeam success", func(t *testing.T) {
		err := repo.AddTeam(teamName, &members)
		require.NoError(t, err)

		team, users, err := repo.GetTeamAndMembers(teamName)
		require.NoError(t, err)
		assert.Equal(t, teamName, team.Name)
		assert.Len(t, *users, len(members))

		truncateAllTables(t, testDB)
	})

	t.Run("AddTeam duplicate returns ErrTeamExists", func(t *testing.T) {
		err := repo.AddTeam(teamName, &members)
		require.NoError(t, err)

		err = repo.AddTeam(teamName, &members)
		require.ErrorIs(t, err, repository.ErrTeamExists)
		truncateAllTables(t, testDB)
	})

	t.Run("GetTeamAndMembers success", func(t *testing.T) {
		err := repo.AddTeam(teamName, &members)
		require.NoError(t, err)

		team, users, err := repo.GetTeamAndMembers(teamName)
		require.NoError(t, err)
		assert.Equal(t, teamName, team.Name)
		assert.Len(t, *users, len(members))
		assert.Equal(t, (*users)[0], members[0])

		truncateAllTables(t, testDB)
	})

	t.Run("GetTeamAndMembers non-existent team returns ErrTeamNotFound", func(t *testing.T) {
		_, _, err := repo.GetTeamAndMembers("nonexistent")
		require.ErrorIs(t, err, repository.ErrTeamNotFound)
		truncateAllTables(t, testDB)
	})
}
