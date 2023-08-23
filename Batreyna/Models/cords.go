package Models

import "gorm.io/gorm"

type LatLng struct {
	gorm.Model
	BranchID uint   `json:"branch_id"`
	Lat      string `json:"lat"`
	Lng      string `json:"lng"`
}
