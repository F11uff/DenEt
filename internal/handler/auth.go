package handler

import (
	"net/http"
	"denet/internal/model"
	"denet/internal/service"
	"denet/internal/http/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	authService service.AuthService
	jwtSecret   string
	logger      *zap.Logger
}

func NewAuthHandler(authService service.AuthService, jwtSecret string, logger *zap.Logger) AuthHandler {
	return &authHandler{
		authService: authService,
		jwtSecret:   jwtSecret,
		logger:      logger,
	}
}

func (h *authHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register request",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to register user",
			zap.String("username", req.Username),
			zap.Error(err),
		)

		switch err {
		case service.ErrUserExists:
			response.WriteError(c, http.StatusConflict, "Username or email already exists")
		default:
			response.WriteError(c, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	token, err := h.authService.GenerateToken(user, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate token",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	h.logger.Info("User registered successfully",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
	)
	response.WriteCreated(c, "User registered successfully", model.AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *authHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	user, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Warn("Failed login attempt",
			zap.String("username", req.Username),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.authService.GenerateToken(user, h.jwtSecret)
	if err != nil {
		h.logger.Error("Failed to generate token",
			zap.String("user_id", user.ID),
			zap.Error(err),
		)
		response.WriteError(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	h.logger.Info("User logged in successfully",
		zap.String("user_id", user.ID),
		zap.String("username", user.Username),
	)
	response.WriteSuccess(c, "Login successful", model.AuthResponse{
		Token: token,
		User:  user,
	})
}