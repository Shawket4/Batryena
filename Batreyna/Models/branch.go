package Models

import "gorm.io/gorm"

type Branch struct {
	gorm.Model
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	LatLng    LatLng    `json:"lat_lng"`
	Inventory Inventory `json:"inventory"`
}
