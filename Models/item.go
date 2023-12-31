package Models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	ParentItemID uint    `json:"parent_item_id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Cost         float64 `json:"cost"`
	IsSold       bool    `json:"is_sold"`
}

type ParentItem struct {
	gorm.Model
	BranchID uint    `json:"branch_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Cost     float64 `json:"cost"`
	Items    []Item  `json:"items"`
}
