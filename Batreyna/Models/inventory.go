package Models

import "gorm.io/gorm"

type Inventory struct {
	gorm.Model
	BranchID uint   `json:"branch_id"`
	Items    []Item `json:"items"`
}
