package Models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ParentItemID uint    `json:"parent_item_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	IsSold       bool    `json:"is_sold"`
}

type ParentItem struct {
	gorm.Model
	BranchID uint   `json:"branch_id"`
	Items    []Item `json:"items"`
}
