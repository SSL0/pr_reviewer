package repository_test

import (
	"context"
	"log"
	"os"
	"pr_reviewer/internal/testutil"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbDSN string

func truncateAllTables(t *testing.T, db *sqlx.DB) {
	_, err := db.Exec("TRUNCATE TABLE users, teams, pull_requests, pull_request_reviewers RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("failed to truncate table: %v", err)
	}
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	tp, err := testutil.StartPostgres(ctx, "../../db/migrations")
	if err != nil {
		log.Fatalf("failed to start postgres: %v", err)
	}
	dbDSN = tp.DSN
	code := m.Run()
	if err := tp.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate postgres: %v", err)
	}

	os.Exit(code)
}
