package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

func TestUserRepository(t *testing.T) {
	testDB, err := repository.NewPostgres(dbDSN)
	require.NoError(t, err)

	userRepo := repository.NewUserRepository(testDB)

	user := model.User{ID: "u19876", Username: "Alice", TeamName: "team161234789", IsActive: true}

	t.Run("SetIsActive success", func(t *testing.T) {
		_, err = testDB.Exec(`INSERT INTO teams(name) VALUES ($1) ON CONFLICT DO NOTHING`, user.TeamName)
		require.NoError(t, err)
		_, err = testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
			user.ID, user.Username, user.TeamName, user.IsActive)

		require.NoError(t, err)
		updated, err := userRepo.SetIsActive(user.ID, false)
		require.NoError(t, err)
		assert.False(t, updated.IsActive)
		assert.Equal(t, user.ID, updated.ID)

		truncateAllTables(t, testDB)
	})

	t.Run("SetIsActive non-existent user returns ErrUserNotFound", func(t *testing.T) {
		_, err := userRepo.SetIsActive("nonexistent", true)
		assert.ErrorIs(t, err, repository.ErrUserNotFound)

		truncateAllTables(t, testDB)
	})

	t.Run("GetAssignedPullRequests success", func(t *testing.T) {
		_, err = testDB.Exec(`INSERT INTO teams(name) VALUES ($1) ON CONFLICT DO NOTHING`, user.TeamName)
		require.NoError(t, err)
		_, err = testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
			user.ID, user.Username, user.TeamName, user.IsActive)

		prID := "pr-1001"
		_, err := testDB.Exec(`
				INSERT INTO pull_requests(id, name, author_id, status, created_at)
				VALUES ($1, $2, $3, $4, NOW())
			`, prID, "Test PR", user.ID, "OPEN")
		require.NoError(t, err)

		_, err = testDB.Exec(`
				INSERT INTO pull_request_reviewers(pull_request_id, reviewer_id)
				VALUES ($1, $2)
			`, prID, user.ID)
		require.NoError(t, err)

		prs, err := userRepo.GetAssignedPullRequests(user.ID)
		require.NoError(t, err)
		assert.Len(t, prs, 1)
		assert.Equal(t, prID, prs[0].ID)

		truncateAllTables(t, testDB)
	})

	t.Run("GetAssignedPullRequests non-existent user returns ErrUserNotFound", func(t *testing.T) {
		_, err := userRepo.GetAssignedPullRequests("nonexistent")
		assert.ErrorIs(t, err, repository.ErrUserNotFound)

		truncateAllTables(t, testDB)
	})
}
