package Models

import (
	"gorm.io/gorm"
	"time"
)

type Shift struct {
	gorm.Model
	BranchID   uint      `json:"branch_id"`
	StartedAt  time.Time `json:"started_at"`
	ClosedAt   time.Time `json:"closed_at"`
	EmployeeID uint      `json:"employee_id"`
	Employee   Employee  `json:"employee" gorm:"-"`
	IsClosed   bool      `json:"is_closed"`
}

type OTP struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  time.Time
	Token      string `json:"token"`
	EmployeeID uint   `json:"employee_id"`
	BranchID   uint   `json:"branch_id"`
}
