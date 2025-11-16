package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

const (
	UniqueViolationCode     = "23505"
	ForeignKeyViolationCode = "23503"
)

func NewPostgres(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil
}
