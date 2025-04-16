package task

import (
	"log"
	"net/http"
	"os"

	"github.com/KiranPawar0/Data-Drive-System/pkg/auth/middleware"
	"github.com/KiranPawar0/Data-Drive-System/pkg/auth/signin"
	"github.com/KiranPawar0/Data-Drive-System/pkg/auth/signup"
	"github.com/KiranPawar0/Data-Drive-System/pkg/fileandfolder"
	"github.com/KiranPawar0/Data-Drive-System/pkg/otpexpiration"
	"github.com/KiranPawar0/Data-Drive-System/routes/getapiroutes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Task(db *gorm.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
		log.Printf("Defaulting to port %s", port)
	}

	apiV1, router := getapiroutes.GetApiRoutes()

	// Define handlers
	apiV1.GET("/auth", func(c *gin.Context) {
		c.String(http.StatusOK, "Task Service Healthy")
	})

	apiV1.POST("/otp-expiration", func(c *gin.Context) {
		otpexpiration.SetOTPExpiration(c, db)
	})

	apiV1.POST("/auth/signup", func(c *gin.Context) {
		signup.SignUpWithEmail(c, db)
	})

	apiV1.POST("/auth/verify-otp", func(c *gin.Context) {
		signup.VerifyOTP(c, db)
	})

	apiV1.POST("/auth/login", func(c *gin.Context) {
		signin.LoginWithEmail(c, db)
	})

	apiV1.POST("/folder/create-folder", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.CreateFolder(c, db)
	})

	apiV1.GET("/folder/get-folder", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.GetFolder(c, db)
	})
	apiV1.PUT("/folder/update-folder", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.UpdateFolder(c, db)
	})
	apiV1.DELETE("/folder/delete-folder", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.DeleteFolder(c, db)
	})
	apiV1.POST("/file/create-file", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.CreateFile(c, db)
	})
	apiV1.GET("/file/get-file", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.GetFile(c, db)
	})
	apiV1.PUT("/file/update-file", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.UpdateFile(c, db)
	})
	apiV1.DELETE("/file/delete-file", middleware.JWTMiddleware(db), func(c *gin.Context) {
		fileandfolder.DeleteFile(c, db)
	})
	// Listen and serve on defined port
	log.Printf("Application started, Listening on Port %s", port)
	router.Run(":" + port)
}
