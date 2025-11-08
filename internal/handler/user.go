package handler

import (
	"denet/internal/http/response"
	"denet/internal/model"
	"denet/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler interface {
	GetUserStatus(c *gin.Context)
	GetLeaderboard(c *gin.Context)
	CompleteTask(c *gin.Context)
	SetReferrer(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService service.UserService, logger *zap.Logger) UserHandler {
	return &userHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *userHandler) GetUserStatus(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.WriteError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	claims, exists := c.Get("user_claims")
	if !exists {
		response.WriteError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwtClaims := claims.(*model.JWTClaims)
	if jwtClaims.UserID != userID && jwtClaims.Username != "admin" {
		response.WriteError(c, http.StatusForbidden, "Access denied")
		return
	}

	userStatus, err := h.userService.GetUserStatus(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user status",
			zap.String("user_id", userID),
			zap.Error(err),
		)

		switch err.Error() {
		case "user not found":
			response.WriteError(c, http.StatusNotFound, "User not found")
		default:
			response.WriteError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	h.logger.Debug("User status retrieved", zap.String("user_id", userID))
	response.WriteSuccess(c, "User status retrieved successfully", userStatus)
}

func (h *userHandler) GetLeaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		response.WriteError(c, http.StatusBadRequest, "Invalid limit parametr")
		return
	}

	if limit > 100 {
		limit = 100
	}

	leaderboard, err := h.userService.GetLeaderboard(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get leaderboard",
			zap.Int("limit", limit),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.logger.Debug("Leaderboard retrieved",
		zap.Int("limit", limit),
		zap.Int("users_count", len(leaderboard)),
	)
	response.WriteSuccess(c, "Leaderboard retrieved successfully", gin.H{
		"leaderboard": leaderboard,
		"limit":       limit,
		"total":       len(leaderboard),
	})
}

func (h *userHandler) CompleteTask(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.WriteError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	claims, exists := c.Get("user_claims")
	if !exists {
		response.WriteError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwtClaims := claims.(*model.JWTClaims)
	if jwtClaims.UserID != userID {
		response.WriteError(c, http.StatusForbidden, "Access denied")
		return
	}

	var req model.CompleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid complete task request",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.CompleteTask(c.Request.Context(), userID, req.TaskID); err != nil {
		h.logger.Error("Failed to complete task",
			zap.String("user_id", userID),
			zap.String("task_id", req.TaskID),
			zap.Error(err),
		)

		switch err.Error() {
		case "user not found":
			response.WriteError(c, http.StatusNotFound, "User not found")
		case "task not found":
			response.WriteError(c, http.StatusBadRequest, "Task not found")
		case "task already completed":
			response.WriteError(c, http.StatusConflict, "Task already completed")
		default:
			response.WriteError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	h.logger.Info("Task completed successfully",
		zap.String("user_id", userID),
		zap.String("task_id", req.TaskID),
	)
	response.WriteSuccess(c, "Task completed successfully", nil)
}

func (h *userHandler) SetReferrer(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		response.WriteError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	claims, exists := c.Get("user_claims")
	if !exists {
		response.WriteError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwtClaims := claims.(*model.JWTClaims)
	if jwtClaims.UserID != userID {
		response.WriteError(c, http.StatusForbidden, "Access denied")
		return
	}

	var req model.SetReferrerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid set referrer request",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.userService.SetReferrer(c.Request.Context(), userID, req.ReferrerID); err != nil {
		h.logger.Error("Failed to set referrer",
			zap.String("user_id", userID),
			zap.String("referrer_id", req.ReferrerID),
			zap.Error(err),
		)

		switch err.Error() {
		case "user not found":
			response.WriteError(c, http.StatusNotFound, "User not found")
		case "referrer already set":
			response.WriteError(c, http.StatusConflict, "Referrer already set")
		case "user cannot refer themselves":
			response.WriteError(c, http.StatusBadRequest, "User cannot refer themselves")
		default:
			response.WriteError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	h.logger.Info("Referrer set successfully",
		zap.String("user_id", userID),
		zap.String("referrer_id", req.ReferrerID),
	)
	response.WriteSuccess(c, "Referrer set successfully", nil)
}
