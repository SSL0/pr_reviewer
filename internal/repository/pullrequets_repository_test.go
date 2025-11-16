package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pr_reviewer/internal/model"
	"pr_reviewer/internal/repository"
)

func TestPullRequestRepository(t *testing.T) {
	testDB, err := repository.NewPostgres(dbDSN)
	require.NoError(t, err)

	prRepo := repository.NewPullRequestRepository(testDB)
	teamName := "team6798678"
	users := []model.User{
		{ID: "u14567", Username: "Alice", TeamName: teamName, IsActive: true},
		{ID: "u23467", Username: "Bob", TeamName: teamName, IsActive: true},
		{ID: "u38043", Username: "Charlie", TeamName: teamName, IsActive: false},
	}
	prID := "pr-113241"

	t.Run("CreatePullRequest success", func(t *testing.T) {
		_, err := testDB.Exec(`INSERT INTO teams(name) VALUES ($1)`, teamName)
		require.NoError(t, err)

		for _, u := range users {
			_, err := testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
				u.ID, u.Username, u.TeamName, u.IsActive)
			require.NoError(t, err)
		}

		prID := prID
		pr, reviewers, err := prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.NoError(t, err)
		assert.Equal(t, prID, pr.ID)
		assert.Equal(t, "Test PR", pr.Name)
		assert.Equal(t, users[0].ID, pr.AuthorID)
		assert.Equal(t, model.PullRequestOpen, pr.Status)
		assert.Len(t, reviewers, 1)

		truncateAllTables(t, testDB)
	})

	t.Run("CreatePullRequest duplicate returns ErrPRExists", func(t *testing.T) {
		_, err := testDB.Exec(`INSERT INTO teams(name) VALUES ($1)`, teamName)
		require.NoError(t, err)

		for _, u := range users {
			_, err := testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
				u.ID, u.Username, u.TeamName, u.IsActive)
			require.NoError(t, err)
		}

		prID := prID
		_, _, err = prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.NoError(t, err)

		_, _, err = prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.ErrorIs(t, err, repository.ErrPRExists)

		truncateAllTables(t, testDB)
	})

	t.Run("CreatePullRequest with non-existent author returns ErrUserNotFound", func(t *testing.T) {
		_, _, err := prRepo.CreatePullRequest(prID, "Test PR", "nonexistent")
		require.ErrorIs(t, err, repository.ErrUserNotFound)

		truncateAllTables(t, testDB)
	})

	t.Run("SetPullRequestStatus success", func(t *testing.T) {
		_, err := testDB.Exec(`INSERT INTO teams(name) VALUES ($1)`, teamName)
		require.NoError(t, err)

		for _, u := range users {
			_, err := testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
				u.ID, u.Username, u.TeamName, u.IsActive)
			require.NoError(t, err)
		}

		prID := prID
		_, _, err = prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.NoError(t, err)

		pr, _, err := prRepo.SetPullRequestStatus(prID, model.PullRequestMerged)
		require.NoError(t, err)
		assert.Equal(t, prID, pr.ID)
		assert.Equal(t, model.PullRequestMerged, pr.Status)

		truncateAllTables(t, testDB)
	})

	t.Run("SetPullRequestStatus non-existent PR returns ErrPRNotFound", func(t *testing.T) {
		_, _, err := prRepo.SetPullRequestStatus("nonexistent", model.PullRequestMerged)
		require.ErrorIs(t, err, repository.ErrPRNotFound)

		truncateAllTables(t, testDB)
	})

	t.Run("ReassignPullRequestReviewer success", func(t *testing.T) {
		_, err := testDB.Exec(`INSERT INTO teams(name) VALUES ($1)`, teamName)
		require.NoError(t, err)

		for _, u := range users {
			_, err := testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
				u.ID, u.Username, u.TeamName, u.IsActive)
			require.NoError(t, err)
		}

		prID := prID

		_, reviewers, err := prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.NoError(t, err)

		_, err = testDB.Exec(`UPDATE users SET is_active = $1 WHERE id = $2`, true, users[2].ID)
		require.NoError(t, err)

		oldReviewer := reviewers[0]

		pr, newReviewers, newReviewerID, err := prRepo.ReassignPullRequestReviewer(prID, oldReviewer)
		require.NoError(t, err)
		assert.Equal(t, prID, pr.ID)
		assert.Contains(t, newReviewers, newReviewerID)

		truncateAllTables(t, testDB)
	})

	t.Run("ReassignPullRequestReviewer without active users in team returns ErrNoCandidate", func(t *testing.T) {
		_, err := testDB.Exec(`INSERT INTO teams(name) VALUES ($1)`, teamName)
		require.NoError(t, err)

		for _, u := range users {
			_, err := testDB.Exec(`INSERT INTO users(id, username, team_name, is_active) VALUES ($1,$2,$3,$4)`,
				u.ID, u.Username, u.TeamName, u.IsActive)
			require.NoError(t, err)
		}

		prID := prID

		_, reviewers, err := prRepo.CreatePullRequest(prID, "Test PR", users[0].ID)
		require.NoError(t, err)
		assert.NotEmpty(t, reviewers)
		assert.Contains(t, reviewers, users[1].ID)

		_, reviewers, _, err = prRepo.ReassignPullRequestReviewer(prID, reviewers[0])
		t.Log(reviewers)
		require.ErrorIs(t, err, repository.ErrNoCanditate)

		truncateAllTables(t, testDB)
	})
}
