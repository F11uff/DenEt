package app

import (
	"denet/config"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(conf *config.Config, logger *zap.Logger) error {
	r := gin.New()

	serverAddr := conf.Server.Host + ":" + conf.Server.Port
	logger.Info("Server starting",
		zap.String("address", serverAddr),
		zap.String("environment", gin.Mode()),
	)

	if err := r.Run(serverAddr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
		return errors.New("Failed to start server")
	}

	return nil
}