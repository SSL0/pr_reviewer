package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"pr_reviewer/internal/config"
	"pr_reviewer/internal/handlers"
	"pr_reviewer/internal/repository"
	"pr_reviewer/internal/service"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to configuration json file")
	flag.Parse()

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	postgres, err := repository.NewPostgres(cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer postgres.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	repo := repository.NewRepository(postgres)
	svc := service.NewService(repo)
	handler := handlers.NewHandler(svc, logger)

	err = handler.RegisterRoutes().Run(cfg.ListeningAddress)
	if err != nil {
		log.Printf("failed to register routes: %v", err)
		return
	}
}
