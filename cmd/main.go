package main

import (
	"log"

	"github.com/kuzmindeniss/itk/internal/config"
	"github.com/kuzmindeniss/itk/internal/db"
	"github.com/kuzmindeniss/itk/internal/db/repository"
	"github.com/kuzmindeniss/itk/internal/handler"
	"github.com/kuzmindeniss/itk/internal/router"
	"github.com/kuzmindeniss/itk/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigrations(cfg); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repository.New(pool)
	walletService := service.NewWalletService(repo)
	walletHandler := handler.NewWalletHandler(walletService)

	r := router.SetupRouter(walletHandler)

	r.Run(":" + cfg.AppPort)
}
