package signin

import (
	"net/http"

	"github.com/KiranPawar0/Data-Drive-System/pkg/auth/config"
	"github.com/KiranPawar0/Data-Drive-System/pkg/helper/jwthelper"
	"github.com/KiranPawar0/Data-Drive-System/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginWithEmail(c *gin.Context, db *gorm.DB) {
	var loginData config.Login
	if err := c.ShouldBindJSON(&loginData); err != nil {
		logrus.WithError(err).Error("Failed to bind login data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	var user models.User
	if err := db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		logrus.WithError(err).Error("User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if user.Email != loginData.Email {
		logrus.WithField("email", loginData.Email).Error("Email mismatch")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		logrus.WithFields(logrus.Fields{
			"email": loginData.Email,
		}).Error("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := jwthelper.GenerateJWTToken(loginData.Email)
	if err != nil {
		logrus.WithError(err).Error("Token generation failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
