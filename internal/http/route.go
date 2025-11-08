package http

import (
	"denet/config"
	"denet/internal/handler"
	"denet/internal/handler/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRoute(authHandler handler.AuthHandler, userHandler handler.UserHandler, conf config.Config, logger *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.CORS())

	public := r.Group("/api")
	{
		public.GET("/health", healthCheck)
		public.POST("/auth/register", authHandler.Register)
		public.POST("/auth/login", authHandler.Login)
		public.GET("/tasks", getTasksList)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(conf.JWT.SecretKey, logger))
	{
		protected.GET("/users/:id/status", userHandler.GetUserStatus)
		protected.GET("/users/leaderboard", userHandler.GetLeaderboard)
		protected.POST("/users/:id/task/complete", userHandler.CompleteTask)
		protected.POST("/users/:id/referrer", userHandler.SetReferrer)
	}

	// 404 handler
	r.NoRoute(notFoundHandler)

	return r
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "OK",
		"service": "user-rewards",
		"version": "1.0.0",
	})
}

func getTasksList(c *gin.Context) {
	tasks := []map[string]interface{}{
		{"id": "1", "name": "referral", "description": "Пригласить друга по реферальному коду", "points": 100},
		{"id": "2", "name": "telegram", "description": "Подписаться на Telegram канал", "points": 50},
		{"id": "3", "name": "twitter", "description": "Подписаться на Twitter", "points": 50},
		{"id": "4", "name": "discord", "description": "Присоединиться к Discord серверу", "points": 75},
		{"id": "5", "name": "profile", "description": "Заполнить профиль", "points": 25},
	}
	c.JSON(200, gin.H{"tasks": tasks})
}

func notFoundHandler(c *gin.Context) {
	c.JSON(404, gin.H{
		"error":   "endpoint not found",
		"message": "check the API documentation for available endpoints",
	})
}
