package testutil

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestPostgres struct {
	Container testcontainers.Container
	DSN       string
	DB        *sqlx.DB
}

func StartPostgres(ctx context.Context, migrationDir string) (*TestPostgres, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s/testdb?sslmode=disable", net.JoinHostPort(host, port.Port()))

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	for range 10 {
		if err := db.Ping(); err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"postgres", driver)
	if err != nil {
		return nil, err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return &TestPostgres{
		Container: container,
		DSN:       dsn,
		DB:        db,
	}, nil
}

func (tp *TestPostgres) Terminate(ctx context.Context) error {
	if err := tp.DB.Close(); err != nil {
		log.Println("failed to close DB:", err)
	}
	return tp.Container.Terminate(ctx)
}
