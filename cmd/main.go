package main

import (
	"denet/config"
	"denet/internal/app"
	

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

	err = app.New(conf, logger)

	if err != nil {
		logger.Fatal("Failed to create app", zap.Error(err))
	}

}
