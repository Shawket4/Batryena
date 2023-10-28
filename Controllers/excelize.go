package Controllers

import (
	"BatrynaBackend/Models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"os"
	"time"
)

func HelloMessage(c *gin.Context) {
	time.Sleep(time.Second * 2)
	c.JSON(http.StatusOK, gin.H{"message": "Shawket2"})
}

func GetBranchTransactionsExcel(c *gin.Context) {
	var input struct {
		BranchID uint `json:"branch_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}

	var branch Models.Branch

	if err := Models.DB.Model(&Models.Branch{}).Preload("Transactions").Where("id = ?", input.BranchID).Find(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	for transactionIndex, transaction := range branch.Transactions {
		var Transaction Models.Transaction
		if err := Models.DB.Model(&Models.Transaction{}).Preload("Items").Where("id = ?", transaction.ID).Find(&Transaction).Error; err != nil {
			ReturnErr(c, err)
		}
		for _, item := range Transaction.Items {
			var Item Models.Item
			if err := Models.DB.Model(&Models.Item{}).Where("id = ?", item.ID).Find(&Item).Error; err != nil {
				ReturnErr(c, err)
			}
			Transaction.ItemsStruct = append(Transaction.ItemsStruct, Item)
		}
		branch.Transactions[transactionIndex] = Transaction
	}

	currentTime := time.Now()
	f := excelize.NewFile()
	fileName := fmt.Sprintf("%s (%s) Transactions_File.xlsx", currentTime, branch.Name)
	filePath := fmt.Sprintf("./ExcelTempFiles/%s", fileName)
	headers := map[string]string{
		"A": "Date",
		"B": "Name",
		"C": "Price",
		"D": "Cost",
		"E": "Profit",
	}
	for cell, value := range headers {
		err := f.SetCellValue("Sheet1", cell+"1", value)
		if err != nil {
			ReturnErr(c, err)
		}
	}
	index := 2
	for transactionIndex, transaction := range branch.Transactions {
		var netProfit float64
		for _, item := range transaction.ItemsStruct {
			f.SetCellValue("Sheet1", fmt.Sprintf("A%v", index), getFormattedDateTime(transaction.CreatedAt))
			f.SetCellValue("Sheet1", fmt.Sprintf("B%v", index), item.Name)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%v", index), item.Price)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%v", index), item.Cost)
			var profit float64 = item.Price - item.Cost
			f.SetCellValue("Sheet1", fmt.Sprintf("E%v", index), profit)
			netProfit += profit
			index++
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", index), "Total Cost")
		f.SetCellValue("Sheet1", fmt.Sprintf("C%v", index), transaction.TotalCost)
		index++
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", index), "Net Profit")
		f.SetCellValue("Sheet1", fmt.Sprintf("C%v", index), netProfit)
		index += 3
		if transactionIndex != len(branch.Transactions)-1 {
			for cell, value := range headers {
				err := f.SetCellValue("Sheet1", fmt.Sprintf("%s%v", cell, index), value)
				if err != nil {
					ReturnErr(c, err)
				}
			}
			index++
		}
	}
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
	}
	byteFile, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", byteFile)
	os.Remove(filePath)
}
