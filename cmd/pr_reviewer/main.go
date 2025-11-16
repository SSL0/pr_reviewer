package main

import (
	"flag"
	"log"
	"pr_reviewer/internal/config"
	"pr_reviewer/internal/handler"
	"pr_reviewer/internal/repository"
	"pr_reviewer/internal/service"
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

	repo := repository.NewRepository(postgres)
	svc := service.NewService(repo)
	handler := handler.NewHandler(svc)

	err = handler.RegisterRoutes().Run(cfg.ListeningAddress)
	if err != nil {
		log.Fatal(err)
	}
}
