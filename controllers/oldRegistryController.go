package controllers

import (
	"net/http"
	"strconv"
	"time"

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

// MigrateOldRegistries handles POST request to copy every not-yet-migrated
// old registry into the registries table. Legacy values are kept verbatim in
// the old_* columns; only t_b, reviewed_at and min_to_mud_date are also
// mapped to their typed counterparts. Idempotent: rows already linked via
// old_registry_id are skipped, so duplicates by t_b are allowed by design.
func MigrateOldRegistries(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	typedUser := user.(models.User)

	var total int64
	if err := initializers.DB.Model(&models.OldRegistry{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count old registries"})
		return
	}

	migratedSubQuery := initializers.DB.Model(&models.Registry{}).
		Select("old_registry_id").
		Where("old_registry_id IS NOT NULL")

	var oldRegistries []models.OldRegistry
	if err := initializers.DB.
		Where("id NOT IN (?)", migratedSubQuery).
		Order("t_b").
		Find(&oldRegistries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch old registries"})
		return
	}

	if len(oldRegistries) == 0 {
		c.JSON(200, gin.H{
			"message":  "Nothing to migrate, all old registries are already migrated",
			"total":    total,
			"migrated": 0,
			"skipped":  total,
		})
		return
	}

	registries := make([]models.Registry, 0, len(oldRegistries))
	for _, old := range oldRegistries {
		oldID := old.ID
		oldTB := old.TB
		tb := int(old.TB)

		registries = append(registries, models.Registry{
			TB:     &tb,
			UserID: &typedUser.ID,
			RegistryDates: models.RegistryDates{
				ReviewedAt: old.SeneSeredilen,
			},
			RegistryMail: models.RegistryMail{
				MinToMudDate: old.SeneHatMinToMud,
			},
			RegistryOldData: models.RegistryOldData{
				OldRegistryID:           &oldID,
				OldTB:                   &oldTB,
				OldMinHat:               old.MinHat,
				OldSeneHatMinToMud:      old.SeneHatMinToMud,
				OldGurujy:               old.Gurujy,
				OldPaychy:               old.Paychy,
				OldSertnamaGurujyPaychy: old.SertnamaGurujyPaychy,
				OldDesga:                old.Desga,
				OldBahaUmumy:            old.BahaUmumy,
				OldMeydanUmumy:          old.MeydanUmumy,
				OldKepResminama:         old.KepResminama,
				OldEmlakPaychy:          old.EmlakPaychy,
				OldBahaPaychy:           old.BahaPaychy,
				OldBaha1m2Paychy:        old.Baha1m2Paychy,
				OldSalgyDesga:           old.SalgyDesga,
				OldSalgyGurujy:          old.SalgyGurujy,
				OldSalgyPaychy:          old.SalgyPaychy,
				OldBashPotr:             old.BashPotr,
				OldSertnamaGurPotr:      old.SertnamaGurPotr,
				OldPotratchyKomek:       old.PotratchyKomek,
				OldShahadatnama:         old.Shahadatnama,
				OldYgtyyarnama:          old.Ygtyyarnama,
				OldPatentPasport:        old.PatentPasport,
				OldSeneBashySongy:       old.SeneBashySongy,
				OldSeneSeredilen:        old.SeneSeredilen,
				OldSeneHasabaAlnan:      old.SeneHasabaAlnan,
				OldWezipeAlanAdam:       old.WezipeAlanAdam,
				OldAdyAlanAdam:          old.AdyAlanAdam,
				OldSeneSanSertnama:      old.SeneSanSertnama,
				OldAdyPaychyAlan:        old.AdyPaychyAlan,
				OldSenePaychyAlan:       old.SenePaychyAlan,
				OldLogin:                old.Login,
			},
		})
	}

	if err := initializers.DB.CreateInBatches(&registries, 100).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to migrate old registries", "details": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message":  "Old registries migrated successfully",
		"total":    total,
		"migrated": len(registries),
		"skipped":  total - int64(len(registries)),
	})
}

// UpdateOldRegistry handles PUT request to update "Alan" fields
func UpdateOldRegistry(c *gin.Context) {
	id := c.Param("id")

	// Find the record first
	var oldRegistry models.OldRegistry
	if err := initializers.DB.Where("t_b = ?", id).First(&oldRegistry).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Old registry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch old registry"})
		return
	}

	// Parse request body
	var requestBody struct {
		WezipeAlanAdam  *string `json:"wezipe_alan_adam"`
		AdyAlanAdam     *string `json:"ady_alan_adam"`
		SeneSanSertnama *string `json:"sene_san_sertnama"`
		AdyPaychyAlan   *string `json:"ady_paychy_alan"`
		SenePaychyAlan  *string `json:"sene_paychy_alan"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields if provided
	updates := make(map[string]interface{})

	if requestBody.WezipeAlanAdam != nil {
		updates["wezipe_alan_adam"] = *requestBody.WezipeAlanAdam
	}

	if requestBody.AdyAlanAdam != nil {
		updates["ady_alan_adam"] = *requestBody.AdyAlanAdam
	}

	if requestBody.SeneSanSertnama != nil {
		// Since the database field is varchar, store as string
		updates["sene_san_sertnama"] = *requestBody.SeneSanSertnama
	}

	if requestBody.AdyPaychyAlan != nil {
		updates["ady_paychy_alan"] = *requestBody.AdyPaychyAlan
	}

	if requestBody.SenePaychyAlan != nil {
		updates["sene_paychy_alan"] = *requestBody.SenePaychyAlan
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Update the record
	if err := initializers.DB.Model(&oldRegistry).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update old registry"})
		return
	}

	// Fetch the updated record
	if err := initializers.DB.Where("t_b = ?", id).First(&oldRegistry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated record"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Old registry updated successfully",
		"data":    oldRegistry,
	})
}
