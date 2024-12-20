package utils

import (
	"github.com/gin-gonic/gin"
)

// PaginatedResponse represents a standard paginated response
type PaginatedResponse struct {
	Data     interface{} `json:"data"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}

// RespondWithPagination sends a standardized paginated response
func RespondWithPagination(c *gin.Context, data interface{}, pagination PaginationParams, total int64) {
	c.JSON(200, PaginatedResponse{
		Data:     data,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		Total:    total,
	})
}
