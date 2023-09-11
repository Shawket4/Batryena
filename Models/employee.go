package Models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name       string `json:"name"`
	Password   string `json:"password"`
	CurrentOTP string `json:"current_otp"`
}
