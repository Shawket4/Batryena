package Models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	Name    string `json:"name"`
	Address string `json:"address"`
	LatLng  LatLng `json:"lat_lng"`
	// Inventory    Inventory     `json:"inventory"`
	ParentItems  []ParentItem  `json:"parent_items"`
	Transactions []Transaction `json:"transactions"`
	HeatMap      HeatMap       `json:"heat_map"`
}

type HeatMap struct {
	gorm.Model
	BranchID  uint    `json:"branch_id"`
	Value     float64 `json:"value"`
	TotalSold float64 `json:"total_sold"`
}
