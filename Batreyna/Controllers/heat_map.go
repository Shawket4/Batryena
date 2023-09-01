package Controllers

import (
	"Batreyna/Models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchBranchesHeatData(c *gin.Context) {
	// var input struct {
	// 	DateFrom string `json:"date_from"`
	// 	DateTo   string `json:"date_to"`
	// }

	// if err := c.ShouldBindJSON(&input); err != nil {
	// 	ReturnErr(c, err)
	// }

	var branches []Models.Branch

	if err := Models.DB.Model(&Models.Branch{}).Preload("LatLng").Preload("HeatMap").Find(&branches).Error; err != nil {
		ReturnErr(c, err)
	}

	for index := range branches {
		var transactions []Models.Transaction
		if err := Models.DB.Model(&Models.Transaction{}).Where("branch_id = ?", branches[index].ID).Find(&transactions).Error; err != nil {
			ReturnErr(c, err)
		}
		for _, transaction := range transactions {
			branches[index].HeatMap.TotalSold += transaction.TotalCost
		}
	}

	branches = CalculateHeatMapValues(branches)

	c.JSON(http.StatusOK, branches)
}
