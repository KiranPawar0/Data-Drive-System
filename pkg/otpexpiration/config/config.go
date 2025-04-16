package config

type OTPExpirationTime struct {
	TimeLimit  int    `json:"timeLimit" validate:"required,min=1"`
	OTPMessage string `json:"otpMessage" validate:"required"`
}
