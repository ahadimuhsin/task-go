package helpers

import (
	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Omit if nil
}

// JSONResponse sends a standardized JSON response in Gin
func JSONResponse(c *gin.Context, statusCode int, success bool, message string, data interface{}) {
	response := StandardResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}