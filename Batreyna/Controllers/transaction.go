package Controllers

import (
	"Batreyna/Models"
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
	if err := c.ShouldBindJSON(input); err != nil {
		ReturnErr(c, err)
	}
	for _, item := range input.Items {
		input.TotalCost += item.Price
	}
	if err := Models.DB.Model(&Models.Transaction{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction Successful"})
}

func UpdateTransaction(c *gin.Context) {
	var input Models.Transaction
	if err := c.ShouldBindJSON(input); err != nil {
		ReturnErr(c, err)
	}
	var transaction Models.Transaction
	if err := Models.DB.Model(&Models.Transaction{}).Where("id = ?", input.ID).Find(&transaction).Error; err != nil {
		ReturnErr(c, err)
	}
	transaction.Items = input.Items

	for _, item := range transaction.Items {
		transaction.TotalCost += item.Price
	}

	if err := Models.DB.Model(&Models.Transaction{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction Successful"})
}

func DeleteTransaction(c *gin.Context) {
	var input struct {
		TransactionID uint `json:"transaction_id"`
	}

	if err := Models.DB.Model(&Models.Transaction{}).Delete("id = ?", input.TransactionID).Error; err != nil {
		ReturnErr(c, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction Deleted"})
}
