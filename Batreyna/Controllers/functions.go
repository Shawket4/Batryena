package Controllers

import (
	"Batreyna/Models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReturnErr(c *gin.Context, err error) {
	log.Println(err.Error())
	c.JSON(http.StatusBadRequest, err)
}

func CalculateHeatMapValues(branches []Models.Branch) []Models.Branch {
	var totalSold float64
	var highestSold float64
	var ratio float64
	for _, branch := range branches {
		if branch.HeatMap.TotalSold > highestSold {
			highestSold = branch.HeatMap.TotalSold
		}
		totalSold += branch.HeatMap.TotalSold
	}
	ratio = highestSold / totalSold
	for index := range branches {
		value := branches[index].HeatMap.TotalSold / totalSold
		value /= ratio
		branches[index].HeatMap.Value = 100 * value
	}
	return branches
}
