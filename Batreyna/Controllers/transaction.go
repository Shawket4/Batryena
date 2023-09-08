package Controllers

import (
	"BatrynaBackend/Models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchTransactions(c *gin.Context) {
	var transactions []Models.Transaction
	if err := Models.DB.Model(&Models.Transaction{}).Find(&transactions).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, transactions)
}

func RegisterTransaction(c *gin.Context) {
	var input Models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	fmt.Println(input)
	var items []Models.Item
	for _, itemID := range input.Items {
		var item Models.Item
		if err := Models.DB.Model(&Models.Item{}).Where("id = ?", itemID.ID).Find(&item).Error; err != nil {
			ReturnErr(c, err)
		}
		items = append(items, item)
	}
	for index := range items {
		items[index].IsSold = true
		input.TotalCost += items[index].Price
	}
	if err := Models.DB.Save(&items).Error; err != nil {
		ReturnErr(c, err)
	}
	if err := Models.DB.Model(&Models.Transaction{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction Successful"})
}

func UpdateTransaction(c *gin.Context) {
	var input Models.Transaction
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	var transaction Models.Transaction
	if err := Models.DB.Model(&Models.Transaction{}).Where("id = ?", input.ID).Find(&transaction).Error; err != nil {
		ReturnErr(c, err)
	}
	var items []Models.Item
	for _, itemID := range input.Items {
		var item Models.Item
		if err := Models.DB.Model(&Models.Item{}).Where("id = ?", itemID.ID).Find(&item).Error; err != nil {
			ReturnErr(c, err)
		}
		items = append(items, item)
	}
	for _, item := range items {
		input.TotalCost += item.Price
	}
	transaction.Items = input.Items

	if err := Models.DB.Model(&Models.Transaction{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction Successful"})
}

func DeleteTransaction(c *gin.Context) {
	var input struct {
		TransactionID uint `json:"transaction_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	if err := Models.DB.Model(&Models.Transaction{}).Delete("id = ?", input.TransactionID).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction Deleted"})
}
