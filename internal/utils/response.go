package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// SuccessResponse sends a successful JSON response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
	})
}

// ErrorResponse sends an error JSON response
func ErrorResponse(c *gin.Context, err *AppError) {
	c.JSON(err.Code, Response{
		Success: false,
		Error: map[string]interface{}{
			"type":    err.Type,
			"message": err.Message,
			"field":   err.Field,
		},
	})
}

// ValidationErrorResponse sends validation errors
func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Error: map[string]interface{}{
			"type":    ErrorTypeValidation,
			"message": "Validation failed",
			"fields":  errors,
		},
	})
}

// SuccessResponseWithMeta sends a successful response with pagination meta
func SuccessResponseWithMeta(c *gin.Context, statusCode int, data interface{}, meta interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}
