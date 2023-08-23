package Models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Items     []Item  `json:"items" gorm:"many2many:transaction_items;"`
	TotalCost float64 `json:"total_cost"`
}
