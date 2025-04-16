package models

import (
	"gorm.io/gorm"
)

// set otp expiration time
type OTPExpirationTime struct {
	gorm.Model
	TimeLimit  int
	OTPMessage string
}

// Models Defined for unverified User
type UnVerifiedUserTemp struct {
	gorm.Model

	Email             string `json:"email" gorm:"not null;unique"`
	Password          string `json:"password" gorm:"not null"`
	UserLoginType     string `json:"userlogintype"`
	OTP               string `json:"otp"`
	OTPExpirationTime int    `json:"otpExpirationTime"`
}

type User struct {
	gorm.Model

	Email             string `json:"email" gorm:"not null;unique"`
	Password          string `json:"password" gorm:"not null"`
	UserLoginType     string `json:"userlogintype"`
	OTPExpirationTime int    `json:"otpExpirationTime"`
}

type Folder struct {
	gorm.Model

	FolderName string   `json:"foldername" gorm:"not null"`                    // Folder name
	ParentID   *uint    `json:"parent_id"`                                     // Nullable, indicates root if nil
	UserID     uint     `json:"user_id" gorm:"not null"`                       // Owner of the folder
	Children   []Folder `gorm:"foreignKey:ParentID" json:"children,omitempty"` // Subfolders
	Files      []File   `gorm:"foreignKey:FolderID" json:"files,omitempty"`    // Files in this folder
}

type File struct {
	gorm.Model

	FileName  string `json:"file_name" gorm:"not null"`
	FolderID  uint   `json:"folder_id" gorm:"not null"`  // foreign key to Folder
	CreatedBy uint   `json:"created_by" gorm:"not null"` // foreign key to User
}
