package Models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	InventoryID uint    `json:"inventory_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Count       int64   `json:"count"`
}
