package Models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Address  string `json:"address"`
	LatLng   LatLng `json:"lat_lng"`
	// Inventory    Inventory     `json:"inventory"`
	ParentItems  []ParentItem  `json:"parent_items"`
	Transactions []Transaction `json:"transactions"`
	HeatMap      HeatMap       `json:"heat_map"`
	TotalSold    float64       `json:"total_sold" gorm:"-"`
	SoldToday    float64       `json:"sold_today" gorm:"-"`
}

type HeatMap struct {
	gorm.Model
	BranchID  uint    `json:"branch_id"`
	Value     float64 `json:"value"`
	TotalSold float64 `json:"total_sold"`
}
