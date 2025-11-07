package main

import (
	"context"
	"denet/config"
	"denet/internal/app"
	pg "denet/internal/store/postgresql"

	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal("failed to init logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting user service", zap.String("Version", "1.0.0"))

	conf, err := config.LoadConfig()
	if err != nil {
		return
	}

	db := pg.NewPostgresDatabase(conf.Database.URL)

	if err := db.Connect(context.Background()); err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	if err := db.RunMigrations(); err != nil {
		logger.Fatal("Failed to run migrations", zap.Error(err))
	}

	err = app.New(conf, logger)

	if err != nil {
		logger.Fatal("Failed to create app", zap.Error(err))
	}

}
