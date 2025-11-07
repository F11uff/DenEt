package app

import (
	"context"
	"denet/config"
	"errors"

	"denet/internal/repository"
	pg "denet/internal/store/postgresql"
	"denet/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func New(conf *config.Config, logger *zap.Logger) error {

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

	uow := repository.NewPostgresUnitOfWork(db)

	userService := service.NewUserService(uow)
	
	//добавить  auth service


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