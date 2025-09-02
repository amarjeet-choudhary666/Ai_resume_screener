package controller

import (
	"errors"
	"net/http"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/database"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"failed to bind json": err.Error()})
			return
		}

		if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields (name, email, password, phone) are required"})
			return
		}

		var existingUser models.User

		if err := database.DB.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := models.HashPassword(newUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to hash password",
			})
			return
		}

		newUser.Password = hashedPassword

		if err := database.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to create user",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"user": gin.H{
			"email":    newUser.Email,
			"password": hashedPassword,
			"name":     newUser.Name,
			"phone":    newUser.Phone,
		},
			"message": "User created"})

	}

}
