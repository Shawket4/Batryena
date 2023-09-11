package Controllers

import (
	"BatrynaBackend/Models"
	"github.com/gin-gonic/gin"
	"net/http"
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

	today := getCurrentFormattedDate()

	for index := range branches {
		var transactionsToday []Models.Transaction

		var transactions []Models.Transaction
		if err := Models.DB.Model(&Models.Transaction{}).Where("branch_id = ? AND DATE(created_at) = ?", branches[index].ID, today).Find(&transactionsToday).Error; err != nil {
			ReturnErr(c, err)
		}

		if err := Models.DB.Model(&Models.Transaction{}).Where("branch_id = ?", branches[index].ID).Find(&transactions).Error; err != nil {
			ReturnErr(c, err)
		}

		for _, transaction := range transactionsToday {
			branches[index].HeatMap.TotalSold += transaction.TotalCost
		}
		for _, transaction := range transactions {
			branches[index].TotalSold += transaction.TotalCost
		}
		branches[index].SoldToday = branches[index].HeatMap.TotalSold
	}

	branches = CalculateHeatMapValues(branches)

	c.JSON(http.StatusOK, branches)
}
