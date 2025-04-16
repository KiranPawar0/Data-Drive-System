package config

type User struct {
	Email         string `json:"email" validate:"required,email"`
	Password      string `json:"password" validate:"required,validpassword"`
	UserLoginType string `json:"userlogintype"`
}

// VerifyOTP Config struct
type VerifyOTP struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,numeric,min=6,max=6"`
}

// Login Config Struct
type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
