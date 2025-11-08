package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteError(c *gin.Context, status int, message string, details ...string) {
	errorResponse := ErrorResponse{
		Error: message,
	}
	if len(details) > 0 {
		errorResponse.Details = details[0]
	}
	c.JSON(status, errorResponse)
}

func WriteSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: message,
		Data:    data,
	})
}

func WriteCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, SuccessResponse{
		Message: message,
		Data:    data,
	})
}
