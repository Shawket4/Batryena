package Models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	BranchID    uint     `json:"branch_id"`
	Items       []ItemID `json:"items" gorm:"many2many:transaction_items;"`
	ItemsStruct []Item   `gorm:"-"`
	TotalCost   float64  `json:"total_cost"`
}

type ItemID struct {
	ID uint `json:"id"`
}
