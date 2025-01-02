package controllers

import (
	"fmt"
	"net/http"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
)

// Recursive function to build the hierarchical structure.
func buildHierarchy(areas []models.Area) []map[string]interface{} {
	var result []map[string]interface{}
	areaMap := make(map[uint][]models.Area)

	// Group areas by ParentID
	for _, area := range areas {
		parentID := uint(0)
		if area.ParentID != nil {
			parentID = *area.ParentID
		}
		areaMap[parentID] = append(areaMap[parentID], area)
	}

	// Recursive function to construct tree
	var constructTree func(parentID uint) []map[string]interface{}
	constructTree = func(parentID uint) []map[string]interface{} {
		var children []map[string]interface{}
		for _, child := range areaMap[parentID] {
			node := map[string]interface{}{
				"value": child.Code,
				"label": child.Name,
			}
			nodeChildren := constructTree(child.Code)
			if len(nodeChildren) > 0 {
				node["children"] = nodeChildren
			}
			children = append(children, node)
		}
		return children
	}

	// Build the hierarchy starting with root (ParentID = 0)
	result = constructTree(0)
	return result
}

// GetAreaHierarchy fetches and returns the hierarchical area data.
func GetAreaHierarchy(c *gin.Context) {
	var areas []models.Area

	// Database connection (replace with your actual DB connection)

	// Fetch all areas
	if err := initializers.DB.Find(&areas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch areas"})
		return
	}

	// Build hierarchical response
	hierarchy := buildHierarchy(areas)

	c.JSON(http.StatusOK, gin.H{"options": hierarchy})
}

// FetchAreaHierarchy fetches the area code and all its parent codes in hierarchical order.
func FetchAreaHierarchy(c *gin.Context) {
	codeParam := c.Param("code")

	// Parse the code parameter into an integer
	var code uint
	if _, err := fmt.Sscanf(codeParam, "%d", &code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code parameter"})
		return
	}

	// Function to recursively fetch parent codes
	var fetchParentCodes func(uint, *[]uint) error
	fetchParentCodes = func(currentCode uint, codes *[]uint) error {
		var area models.Area
		if err := initializers.DB.Where("code = ?", currentCode).First(&area).Error; err != nil {
			return err
		}

		// Prepend the current code
		*codes = append([]uint{area.Code}, *codes...)

		// If there is a parent, fetch recursively
		if area.ParentID != nil {
			return fetchParentCodes(*area.ParentID, codes)
		}
		return nil
	}

	// Fetch the hierarchy
	var codes []uint
	if err := fetchParentCodes(code, &codes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hierarchy"})
		return
	}

	// Respond with the codes
	c.JSON(http.StatusOK, gin.H{"codes": codes})
}
