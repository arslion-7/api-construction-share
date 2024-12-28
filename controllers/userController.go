package controllers

import (
	"net/http"

	"github.com/arslion-7/api-construction-share/initializers"
	"github.com/arslion-7/api-construction-share/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context) {
	show := c.Query("show")
	var users []models.User

	query := initializers.DB.Order("id")

	if show == "with-deleted" {
		if err := query.Unscoped().Find(&users).Error; err != nil {
			c.AbortWithStatus(404)
			return
		} else {
			c.JSON(200, users)
			return
		}
	}
	if err := query.Find(&users).Error; err != nil {
		c.AbortWithStatus(404)
		return
	} else {
		c.JSON(200, users)
	}
}

func GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user models.User
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
	} else {
		c.JSON(200, user)
	}
}

func CreateUser(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	// Create user record in the database
	newUser := models.User{
		Email:       input.Email,
		Password:    string(hashedPassword),
		FullName:    input.FullName,
		PhoneNumber: input.PhoneNumber,
		Role:        input.Role,
	}
	if err := initializers.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

type UserUpdateRequest struct {
	Email       string  `json:"email" binding:"required"`
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func UpdateUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	var input UserUpdateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
	// 	return
	// }

	user.Email = input.Email
	user.FullName = input.FullName
	user.PhoneNumber = input.PhoneNumber
	// user.Password = string(hashedPassword)

	// initializers.DB.Omit("Password").Save(&user)
	initializers.DB.Omit("Password").Save(&user)
	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")
	if err := initializers.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	initializers.DB.Delete(&user)
	c.JSON(200, "User deleted successfully")
}

func RestoreUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")
	if err := initializers.DB.Unscoped().Where("id = ?", id).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	initializers.DB.Model(&user).Unscoped().Select("DeletedAt").Updates(map[string]interface{}{"DeletedAt": nil})
	c.JSON(200, "User restored successfully")
}

// type CreateUserInput struct {
// 	Email    string  `json:"email" binding:"required"`
// 	FullName *string `json:"full_name" binding:"required"`
// 	Password string  `json:"password" binding:"required"`
// }

// func CreateUser(c *gin.Context) {
// 	var input CreateUserInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}
// 	// Create user record in the database
// 	newUser := models.User{
// 		Email:    input.Email,
// 		Password: string(hashedPassword),
// 		FullName: input.FullName,
// 	}
// 	if err := initializers.DB.Create(&newUser).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	c.JSON(http.StatusCreated, newUser)
// }
