package signup

import (
	"net/http"
	"time"

	"github.com/KiranPawar0/Data-Drive-System/pkg/auth/config"
	"github.com/KiranPawar0/Data-Drive-System/pkg/helper/otphelper"
	"github.com/KiranPawar0/Data-Drive-System/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUpWithEmail(c *gin.Context, db *gorm.DB) {
	var signUpData config.User

	if err := c.ShouldBindJSON(&signUpData); err != nil {
		logrus.WithError(err).Error("Failed to bind signup data")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingVerifiedUser models.User
	if err := db.Where("email = ?", signUpData.Email).First(&existingVerifiedUser).Error; err == nil {
		logrus.WithField("email", signUpData.Email).Info("Email already exists and is verified")
		c.JSON(http.StatusOK, gin.H{"message": "Email already exists and is verified. Please login."})
		return
	}

	var existingUnverifiedUser models.UnVerifiedUserTemp
	unverifiedCheck := db.Where("email = ?", signUpData.Email).First(&existingUnverifiedUser)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpData.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.WithError(err).Error("Password hashing failed")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to hash password"})
		return
	}

	otp := otphelper.GenerateOTP()
	_, timeLimit, err := otphelper.GetOTPTimeLimitAndMessage(c, db)
	if err != nil {
		logrus.WithError(err).Error("Failed to get OTP time limit")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to get OTP time limit"})
		return
	}

	expiryTime := int(time.Now().Add(time.Duration(timeLimit) * time.Minute).Unix())

	if unverifiedCheck.Error == nil {
		err = otphelper.SendOTPEmail(signUpData.Email, otp, "Your OTP for Dating App Email Verification is:", "authentication")
		if err != nil {
			logrus.WithError(err).Error("Failed to send OTP email")
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to send OTP email"})
			return
		}

		updateData := map[string]interface{}{
			"otp":                 otp,
			"otp_expiration_time": expiryTime,
			"password":            string(hashedPassword),
		}

		if err := db.Model(&models.UnVerifiedUserTemp{}).Where("email = ?", signUpData.Email).Updates(updateData).Error; err != nil {
			logrus.WithError(err).Error("Failed to update OTP")
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to update OTP"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User already exists. OTP sent to registered email."})
		return
	}

	userinfo := models.UnVerifiedUserTemp{
		Email:             signUpData.Email,
		Password:          string(hashedPassword),
		UserLoginType:     signUpData.UserLoginType,
		OTP:               otp,
		OTPExpirationTime: expiryTime,
	}

	err = otphelper.SendOTPEmail(signUpData.Email, otp, "Your OTP for Data Drive System Email Verification is:", "authentication")
	if err != nil {
		logrus.WithError(err).Error("Failed to send OTP email")
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to send OTP email"})
		return
	}

	if err := db.Create(&userinfo).Error; err != nil {
		logrus.WithError(err).Error("Failed to insert user info")
		c.JSON(http.StatusBadRequest, gin.H{"errors": "Failed to insert user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully. OTP sent to registered email."})
}

func VerifyOTP(c *gin.Context, db *gorm.DB) {
	var verifyOTPData config.VerifyOTP
	if err := c.ShouldBindJSON(&verifyOTPData); err != nil {
		logrus.WithError(err).Error("Failed to bind OTP data")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.UnVerifiedUserTemp
	if err := db.Where("email = ?", verifyOTPData.Email).First(&existingUser).Error; err != nil {
		logrus.WithField("email", verifyOTPData.Email).Error("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if existingUser.OTP != verifyOTPData.OTP || time.Now().Unix() > int64(existingUser.OTPExpirationTime) {
		if time.Now().Unix() > int64(existingUser.OTPExpirationTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		}
		return
	}

	verifiedUser := models.User{
		Email:         existingUser.Email,
		Password:      existingUser.Password,
		UserLoginType: existingUser.UserLoginType,
	}

	if err := db.Create(&verifiedUser).Error; err != nil {
		logrus.WithError(err).Error("Failed to insert verified user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user"})
		return
	}

	if err := db.Where("email = ?", verifyOTPData.Email).Delete(&models.UnVerifiedUserTemp{}).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete unverified user entry")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clean up unverified data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verified successfully"})
}
