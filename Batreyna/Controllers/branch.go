package Controllers

import (
	"Batreyna/Models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchBranches(c *gin.Context) {
	var branches []Models.Branch
	if err := Models.DB.Model(Models.Branch{}).Preload("LatLng").Preload("ParentItems").Preload("Transactions").Preload("HeatMap").Find(&branches).Error; err != nil {
		ReturnErr(c, err)
	}
	for branchIndex := range branches {
		// var inventory Models.Inventory
		var transactions []Models.Transaction
		// if err := Models.DB.Model(&Models.Inventory{}).Preload("Items").Where("id = ?", branches[index].Inventory.ID).Find(&inventory).Error; err != nil {
		// 	ReturnErr(c, err)
		// }
		// branches[index].Inventory = inventory
		for parentItemIndex := range branches[branchIndex].ParentItems {
			var parentItem Models.ParentItem
			if err := Models.DB.Model(&Models.ParentItem{}).Preload("Items").Where("id = ?", branches[branchIndex].ParentItems[parentItemIndex].ID).Find(&parentItem).Error; err != nil {
				ReturnErr(c, err)
			}
			branches[branchIndex].ParentItems[parentItemIndex] = parentItem
		}
		if err := Models.DB.Model(&Models.Transaction{}).Preload("Items").Where("branch_id = ?", branches[branchIndex].ID).Find(&transactions).Error; err != nil {
			ReturnErr(c, err)
		}
		branches[branchIndex].Transactions = transactions
	}
	c.JSON(http.StatusOK, branches)
}

func RegisterBranch(c *gin.Context) {
	var input Models.Branch

	if err := c.BindJSON(&input); err != nil {
		ReturnErr(c, err)
		return
	}
	// input.Inventory = Models.Inventory{Items: []Models.Item{}}
	input.HeatMap = Models.HeatMap{}
	input.HeatMap.Value = 1
	if err := Models.DB.Model(&Models.Branch{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Registered"})
}

func UpdateBranch(c *gin.Context) {
	var input Models.Branch
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}

	var branch Models.Branch

	if err := Models.DB.Model(&Models.Branch{}).Where("id = ?", input.ID).Preload("ParentItems").Find(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	if err := Models.DB.Model(&Models.ParentItem{}).Delete(&branch.ParentItems).Error; err != nil {
		ReturnErr(c, err)
	}

	branch.Name = input.Name
	branch.LatLng.Lat = input.LatLng.Lat
	branch.LatLng.Lng = input.LatLng.Lng
	branch.ParentItems = input.ParentItems
	branch.Address = input.Address

	if err := Models.DB.Save(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Updated"})
}

func DeleteBranch(c *gin.Context) {
	var input struct {
		BranchID uint `json:"branch_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	if err := Models.DB.Model(&Models.Branch{}).Delete("id = ?", input.BranchID).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Deleted"})
}
