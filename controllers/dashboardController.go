package controllers

import (
	"net/http"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
)

type UserRegistryStats struct {
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	Count     int64  `json:"count"`
}

type DashboardStats struct {
	TotalRegistries     int64               `json:"total_registries"`
	TotalUsers          int64               `json:"total_users"`
	RegistriesByUser    []UserRegistryStats `json:"registries_by_user"`
	RegistriesThisMonth int64               `json:"registries_this_month"`
	RegistriesThisYear  int64               `json:"registries_this_year"`
}

// GetDashboardStats handles GET request for dashboard statistics
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// Get total registries count
	if err := initializers.DB.Model(&models.Registry{}).Count(&stats.TotalRegistries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total registries count"})
		return
	}

	// Get total users count
	if err := initializers.DB.Model(&models.User{}).Count(&stats.TotalUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total users count"})
		return
	}

	// Get registries count by user
	var userStats []UserRegistryStats
	if err := initializers.DB.Table("registries").
		Select("registries.user_id, COALESCE(users.full_name, 'Unknown User') as user_name, users.email as user_email, COUNT(*) as count").
		Joins("LEFT JOIN users ON registries.user_id = users.id").
		Where("registries.deleted_at IS NULL").
		Group("registries.user_id, users.full_name, users.email").
		Order("count DESC").
		Find(&userStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get registries by user"})
		return
	}
	stats.RegistriesByUser = userStats

	// Get registries created this month
	if err := initializers.DB.Model(&models.Registry{}).
		Where("created_at >= DATE_TRUNC('month', CURRENT_DATE)").
		Count(&stats.RegistriesThisMonth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get this month's registries count"})
		return
	}

	// Get registries created this year
	if err := initializers.DB.Model(&models.Registry{}).
		Where("created_at >= DATE_TRUNC('year', CURRENT_DATE)").
		Count(&stats.RegistriesThisYear).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get this year's registries count"})
		return
	}

	c.JSON(200, stats)
}
