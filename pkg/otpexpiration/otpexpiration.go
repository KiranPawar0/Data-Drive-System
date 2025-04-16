package otpexpiration

import (
	"net/http"

	"github.com/KiranPawar0/Data-Drive-System/pkg/models"
	"github.com/KiranPawar0/Data-Drive-System/pkg/otpexpiration/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// SetOTPExpiration updates the OTP expiration time
func SetOTPExpiration(c *gin.Context, db *gorm.DB) {
	var otpConfigData config.OTPExpirationTime

	// Bind JSON
	if err := c.ShouldBindJSON(&otpConfigData); err != nil {
		logrus.WithError(err).Warn("Failed to bind JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		return
	}

	// Prepare model
	otpData := models.OTPExpirationTime{
		TimeLimit:  otpConfigData.TimeLimit,
		OTPMessage: otpConfigData.OTPMessage,
	}

	// Save to DB
	if err := db.Create(&otpData).Error; err != nil {
		logrus.WithError(err).Error("Failed to save OTP expiration time")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP expiration time"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"TimeLimit":  otpData.TimeLimit,
		"OTPMessage": otpData.OTPMessage,
	}).Info("OTP expiration time saved successfully")

	c.JSON(http.StatusOK, gin.H{"message": "OTP expiration time updated successfully"})
}

// Get latest OTP expiration time
func GetOTPExpiration(c *gin.Context, db *gorm.DB) {
	var otpData models.OTPExpirationTime

	if err := db.Order("created_at desc").First(&otpData).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logrus.Warn("No OTP expiration time found")
			c.JSON(http.StatusNotFound, gin.H{"error": "No OTP expiration time found"})
		} else {
			logrus.WithError(err).Error("Failed to retrieve OTP expiration time")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve OTP expiration time"})
		}
		return
	}

	logrus.WithField("TimeLimit", otpData.TimeLimit).Info("OTP expiration time retrieved")
	c.JSON(http.StatusOK, gin.H{
		"timelimit":  otpData.TimeLimit,
		"otpMessage": otpData.OTPMessage,
	})
}
