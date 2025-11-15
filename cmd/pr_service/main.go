package main

import (
	"flag"
	"log"
	"pr_service/internal/config"
	"pr_service/internal/handler"
	"pr_service/internal/repository"
	"pr_service/internal/service"
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
	service := service.NewService(repo)
	handler := handler.NewHandler(service)

	err = handler.RegisterRoutes().Run(cfg.ListeningAddress)
	if err != nil {
		log.Fatal(err)
	}
}
