#!/bin/sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3
migrate -database "$POSTGRES_URL?sslmode=disable" -path db/migrations/ up
