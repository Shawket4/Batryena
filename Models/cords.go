package Models

import "gorm.io/gorm"

type LatLng struct {
	gorm.Model
	BranchID uint    `json:"branch_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
}
