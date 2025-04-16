package otphelper

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/KiranPawar0/Data-Drive-System/pkg/models"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

// SendOTPEmail sends an OTP to the specified email address.
func SendOTPEmail(email, otp, message, otpUsage string) error {
	username := strings.ToUpper(string(email[0])) + strings.ToLower(strings.Split(email, "@")[0][1:])

	var subject string
	if otpUsage == "authentication" {
		subject = "Data Drive System OTP For Email Verification"
	}

	// Set up email configuration
	m := gomail.NewMessage()
	m.SetHeader("From", "Data Drive System <kiran.pawar@bharatverified.in>")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	var htmlBody string
	if otpUsage == "authentication" {
		htmlBody = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="utf-8">
			<title>OTP Verification</title>
			<style>
				body { font-family: Arial, sans-serif; background-color: #f5f5f5; }
				.card { max-width: 600px; margin: auto; background: white; padding: 20px; border-radius: 10px; }
				h1 { color: #333; text-align: center; }
				p { text-align: center; font-size: 16px; color: #666; }
				.otp { font-size: 24px; font-weight: bold; color: #000; }
			</style>
		</head>
		<body>
			<div class="card">
				<h1>Hello, ` + username + `</h1>
				<p>` + message + `</p>
				<p class="otp">` + otp + `</p>
				<p>Please use this OTP to complete your ` + otpUsage + `.</p>
				<p>If you didn't request this OTP, please ignore this email.</p>
			</div>
		</body>
		</html>`
	}

	m.SetBody("text/html", htmlBody)

	// Email sending configuration
	d := gomail.NewDialer("smtp.zoho.com", 587, "YOUR_EMAIL", "YOUR_PASS")
	// Enable TLS (recommended for Gmail)
	// d.TLSConfig = &tls.Config{
	// 	InsecureSkipVerify: false,
	// }

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Error sending OTP email:", err)
		return err
	}
	return nil
}

// func to generate OTP
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())

	otp := ""
	for i := 0; i < 6; i++ {
		otp += fmt.Sprint(rand.Intn(10))
	}

	return otp
}

func GetOTPTimeLimitAndMessage(c *gin.Context, db *gorm.DB) (string, int, error) {
	var otpData models.OTPExpirationTime

	// Get the latest entry (assuming latest by ID or created_at)
	err := db.Order("id desc").First(&otpData).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"errors": "No OTP expiration time found"})
			return "", 0, err
		}
		c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to retrieve OTP expiration time"})
		return "", 0, err
	}

	return otpData.OTPMessage, otpData.TimeLimit, nil
}
