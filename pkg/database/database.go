package database

import (
	"fmt"
	"time"

	"github.com/KiranPawar0/Data-Drive-System/pkg/config"
	"github.com/KiranPawar0/Data-Drive-System/pkg/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	var err error

	// Open database connection
	cfg, err := config.Env()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DatabaseName)
	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Set database connection settings
	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)                 // Maximum number of idle connections in the connection pool
	sqlDB.SetMaxOpenConns(100)                // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	// Create database tables
	err = dbConn.AutoMigrate(&models.UnVerifiedUserTemp{}, models.OTPExpirationTime{}, models.User{}, models.Folder{}, models.File{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate User table: %v", err)
	}

	// err = dbConn.AutoMigrate(&models.Tenant{}, &models.User{}, &models.Attendance{}, &models.GPS{}, &models.OutOfRange{}, &models.Message{}, &models.LoginHours{}, &models.Department{}, &models.RoleStatusCode{}, &models.Department{}, &models.LeaveStatus{}, &models.LeaveBalance{}, &models.LeaveCategories{}, &models.OfficeLocation{}, &models.Projects{}, &models.ProjectStatus{}, &models.ProjectParticipent{}, &models.TTSUsersRole{}, &models.Tasks{}, &models.TaskStatus{}, &models.OfficeStatus{}, &models.TaskDuration{}, &models.Outdoor{})
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to auto migrate other database tables: %v", err)
	// }

	return dbConn, nil
}
