package controllers

import (
	"net/http"
	"strconv"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetOldRegistries handles GET request for old registries with pagination
func GetOldRegistries(c *gin.Context) {
	var oldRegistries []models.OldRegistry
	var total int64

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	orderBy := c.DefaultQuery("orderBy", "t_b")
	orderDir := c.DefaultQuery("orderDir", "asc")

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	query := initializers.DB.Model(&models.OldRegistry{})

	// Add search functionality
	if search != "" {
		if tb, err := strconv.Atoi(search); err == nil {
			// Search by t_b if the search term can be converted to an integer
			query = query.Where("t_b = ?", tb)
		} else {
			// Search by text fields if the search term is not a number
			searchQuery := "%" + search + "%"
			query = query.Where(
				initializers.DB.Where("min_hat ILIKE ?", searchQuery).
					Or("gurujy ILIKE ?", searchQuery).
					Or("paychy ILIKE ?", searchQuery).
					Or("desga ILIKE ?", searchQuery).
					Or("salgy_desga ILIKE ?", searchQuery).
					Or("salgy_gurujy ILIKE ?", searchQuery).
					Or("salgy_paychy ILIKE ?", searchQuery).
					Or("login ILIKE ?", searchQuery),
			)
		}
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count records"})
		return
	}

	// Build order clause
	orderClause := orderBy + " " + orderDir
	if orderDir != "asc" && orderDir != "desc" {
		orderClause = "t_b asc" // Default fallback
	}

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order(orderClause).Find(&oldRegistries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch old registries"})
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrev := page > 1

	c.JSON(200, gin.H{
		"data": oldRegistries,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
			"hasNext":    hasNext,
			"hasPrev":    hasPrev,
		},
	})
}

// GetOldRegistry handles GET request for a single old registry
func GetOldRegistry(c *gin.Context) {
	id := c.Param("id")
	var oldRegistry models.OldRegistry

	if err := initializers.DB.Where("t_b = ?", id).First(&oldRegistry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Old registry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch old registry"})
		return
	}

	c.JSON(200, oldRegistry)
}
