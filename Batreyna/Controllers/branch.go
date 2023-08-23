package Controllers

import (
	"Batreyna/Models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchBranches(c *gin.Context) {
	var branches []Models.Branch
	if err := Models.DB.Model(Models.Branch{}).Preload("LatLng").Preload("Inventory").Find(&branches).Error; err != nil {
		ReturnErr(c, err)
	}
	for index := range branches {
		var inventory Models.Inventory
		if err := Models.DB.Model(&Models.Inventory{}).Preload("Items").Where("id = ?", branches[index].Inventory.ID).Find(&inventory).Error; err != nil {
			ReturnErr(c, err)
		}
		branches[index].Inventory = inventory
	}
	c.JSON(http.StatusOK, branches)
}

func RegisterBranch(c *gin.Context) {
	var input Models.Branch

	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}

	if err := Models.DB.Model(&Models.Branch{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Registered"})
}

func UpdateBranch(c *gin.Context) {
	var input Models.Branch
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}

	var branch Models.Branch

	if err := Models.DB.Model(&Models.Branch{}).Where("id = ?", input.ID).Find(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	branch.Name = input.Name
	branch.LatLng = input.LatLng
	branch.Inventory = input.Inventory

	if err := Models.DB.Save(&branch.LatLng).Error; err != nil {
		ReturnErr(c, err)
	}

	if err := Models.DB.Save(&branch.Inventory).Error; err != nil {
		ReturnErr(c, err)
	}

	if err := Models.DB.Save(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Registered"})
}

func DeleteBranch(c *gin.Context) {
	var input struct {
		BranchID uint `json:"branch_id"`
	}

	if err := Models.DB.Model(&Models.Branch{}).Delete("id = ?", input.BranchID).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Deleted"})
}
