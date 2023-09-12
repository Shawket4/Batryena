package Controllers

import (
	"BatrynaBackend/Models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FetchBranchList(c *gin.Context) {
	var branches []Models.Branch
	if err := Models.DB.Model(Models.Branch{}).Find(&branches).Error; err != nil {
		ReturnErr(c, err)
	}
	for index := range branches {
		branches[index].TotalSold = 0
		branches[index].SoldToday = 0
		branches[index].Password = ""
	}
	c.JSON(http.StatusOK, branches)
}

func FetchBranches(c *gin.Context) {
	var branches []Models.Branch
	if err := Models.DB.Model(Models.Branch{}).Preload("LatLng").Preload("ParentItems").Preload("Transactions").Preload("HeatMap").Preload("Shifts").Find(&branches).Error; err != nil {
		ReturnErr(c, err)
	}
	for branchIndex := range branches {
		if len(branches[branchIndex].Shifts) != 0 {
			var currentShift Models.Shift
			if err := Models.DB.Model(&Models.Shift{}).Where("branch_id = ?", branches[branchIndex].ID).Last(&currentShift).Error; err != nil {
				ReturnErr(c, err)
			}
			var currentEmployee Models.Employee
			if err := Models.DB.Model(&Models.Employee{}).Where("id = ?", currentShift.EmployeeID).Find(&currentEmployee).Error; err != nil {
				ReturnErr(c, err)
			}
			currentShift.Employee = currentEmployee
			branches[branchIndex].CurrentShift = currentShift
		}
		for shiftIndex := range branches[branchIndex].Shifts {
			var Employee Models.Employee
			if err := Models.DB.Model(&Models.Employee{}).Where("id = ?", branches[branchIndex].Shifts[shiftIndex].EmployeeID).Find(&Employee).Error; err != nil {
				ReturnErr(c, err)
			}
			branches[branchIndex].Shifts[shiftIndex].Employee = Employee
		}

		// var inventory Models.Inventory
		var transactions []Models.Transaction
		// if err := Models.DB.Model(&Models.Inventory{}).Preload("Items").Where("id = ?", branches[index].Inventory.ID).Find(&inventory).Error; err != nil {
		// 	ReturnErr(c, err)
		// }
		// branches[index].Inventory = inventory
		for parentItemIndex := range branches[branchIndex].ParentItems {
			var ParentItem Models.ParentItem
			var SubItems []Models.Item
			if err := Models.DB.Model(&Models.ParentItem{}).Preload("Items").Where("id = ?", branches[branchIndex].ParentItems[parentItemIndex].ID).Find(&ParentItem).Error; err != nil {
				ReturnErr(c, err)
			}
			for _, subItem := range ParentItem.Items {
				if !subItem.IsSold {
					SubItems = append(SubItems, subItem)
				}
			}
			ParentItem.Items = SubItems
			branches[branchIndex].ParentItems[parentItemIndex] = ParentItem
		}
		if err := Models.DB.Model(&Models.Transaction{}).Preload("Items").Where("branch_id = ?", branches[branchIndex].ID).Find(&transactions).Error; err != nil {
			ReturnErr(c, err)
		}
		for transactionIndex, transaction := range transactions {
			var ItemStructs []Models.Item
			for _, itemID := range transaction.Items {
				var Item Models.Item
				if err := Models.DB.Model(&Models.Item{}).Where("id = ?", itemID).Find(&Item).Error; err != nil {
					ReturnErr(c, err)
				}
				ItemStructs = append(ItemStructs, Item)
			}
			transactions[transactionIndex].ItemsStruct = ItemStructs
		}
		branches[branchIndex].Transactions = transactions
		var transactionsToday []Models.Transaction
		today := getCurrentFormattedDate()
		if err := Models.DB.Model(&Models.Transaction{}).Where("branch_id = ? AND DATE(created_at) = ?", branches[branchIndex].ID, today).Find(&transactionsToday).Error; err != nil {
			ReturnErr(c, err)
		}

		for _, transaction := range transactionsToday {
			branches[branchIndex].HeatMap.TotalSold += transaction.TotalCost
		}
		for _, transaction := range transactions {
			branches[branchIndex].TotalSold += transaction.TotalCost
		}
		branches[branchIndex].SoldToday = branches[branchIndex].HeatMap.TotalSold
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
	var user Models.User
	user.Permission = 1
	user.Username = input.Name
	user.Password = input.Password
	input.Password = ""
	if err := Models.DB.Model(&Models.Branch{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
		return
	}
	var branch Models.Branch
	if err := Models.DB.Model(&Models.Branch{}).Where("name = ?", input.Name).Find(&branch).Error; err != nil {
		ReturnErr(c, err)
		return
	}
	user.BranchID = branch.ID
	_, err := user.SaveUser()
	if err != nil {
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

	branch.Name = input.Name
	branch.LatLng.Lat = input.LatLng.Lat
	branch.LatLng.Lng = input.LatLng.Lng
	branch.Address = input.Address
	var inputItems []Models.Item
	var branchItems []Models.Item
	for _, parentItem := range input.ParentItems {
		for _, subItem := range parentItem.Items {
			inputItems = append(inputItems, subItem)
		}
	}
	for _, parentItem := range branch.ParentItems {
		var subItems []Models.Item
		if err := Models.DB.Model(&Models.Item{}).Where("parent_item_id = ?", parentItem.ID).Find(&subItems).Error; err != nil {
			ReturnErr(c, err)
		}
		for _, subItem := range subItems {
			branchItems = append(branchItems, subItem)
		}
	}

	var deletedItems []Models.Item

	for _, subItem := range branchItems {
		var isFound bool
		for _, subItemInput := range inputItems {
			if subItemInput.ID == subItem.ID {
				isFound = true
			}
		}
		if !isFound {
			deletedItems = append(deletedItems, subItem)
		}
	}
	fmt.Println(deletedItems)

	var deletedParentItems []Models.ParentItem
	for _, parentItem := range branch.ParentItems {
		var isFound bool
		for _, parentItemInput := range input.ParentItems {
			if parentItem.ID == parentItemInput.ID {
				isFound = true
			}
		}
		if !isFound {
			deletedParentItems = append(deletedParentItems, parentItem)
		}
	}
	if len(deletedItems) != 0 {
		if err := Models.DB.Model(&Models.Item{}).Unscoped().Delete(&deletedItems).Error; err != nil {
			return
		}
	}
	if len(deletedParentItems) != 0 {
		if err := Models.DB.Model(&Models.ParentItem{}).Unscoped().Delete(&deletedParentItems).Error; err != nil {
			return
		}
	}
	//var deletedItems []Models.Item
	//for _, subItem := range branchItems {
	//	var isFound bool = false
	//	for _, subItemKept := range deletedItems {
	//		if subItem.ID == subItemKept.ID {
	//
	//		}
	//	}
	//}
	//var deletedParentItems []Models.Item
	//var items []Models.Item
	//for _, parentItem := range branch.ParentItems {
	//	var ParentItem Models.ParentItem
	//	if err := Models.DB.Model(&Models.ParentItem{}).Preload("Items").Where("id = ?", parentItem.ID).Find(&ParentItem).Error; err != nil {
	//		ReturnErr(c, err)
	//	}
	//	for _, item := range ParentItem.Items {
	//		items = append(items, item)
	//	}
	//
	//}
	//if len(branch.ParentItems) != 0 {
	//	if err := Models.DB.Model(&Models.ParentItem{}).Unscoped().Delete(&branch.ParentItems).Error; err != nil {
	//		return
	//	}
	//}
	//
	//if len(items) != 0 {
	//	if err := Models.DB.Model(&Models.Item{}).Unscoped().Delete(&items).Error; err != nil {
	//		return
	//	}
	//}

	branch.ParentItems = input.ParentItems

	if err := Models.DB.Save(&branch).Error; err != nil {
		ReturnErr(c, err)
	}

	//if err := Models.DB.Save(&branch.ParentItems).Error; err != nil {
	//	ReturnErr(c, err)
	//}

	c.JSON(http.StatusOK, gin.H{"message": "Branch Updated"})
}

//func UpdateBranch(c *gin.Context) {
//	var input Models.Branch
//	if err := c.ShouldBindJSON(&input); err != nil {
//		ReturnErr(c, err)
//	}
//
//	var branch Models.Branch
//
//	if err := Models.DB.Model(&Models.Branch{}).Where("id = ?", input.ID).Preload("ParentItems").Find(&branch).Error; err != nil {
//		ReturnErr(c, err)
//	}
//
//	branch.Name = input.Name
//	branch.LatLng.Lat = input.LatLng.Lat
//	branch.LatLng.Lng = input.LatLng.Lng
//	branch.Address = input.Address
//	var items []Models.Item
//	for _, parentItem := range branch.ParentItems {
//		var ParentItem Models.ParentItem
//		if err := Models.DB.Model(&Models.ParentItem{}).Preload("Items").Where("id = ?", parentItem.ID).Find(&ParentItem).Error; err != nil {
//			ReturnErr(c, err)
//		}
//		for _, item := range ParentItem.Items {
//			items = append(items, item)
//		}
//
//	}
//	if len(branch.ParentItems) != 0 {
//		if err := Models.DB.Model(&Models.ParentItem{}).Unscoped().Delete(&branch.ParentItems).Error; err != nil {
//			return
//		}
//	}
//
//	if len(items) != 0 {
//		if err := Models.DB.Model(&Models.Item{}).Unscoped().Delete(&items).Error; err != nil {
//			return
//		}
//	}
//
//	branch.ParentItems = input.ParentItems
//
//	fmt.Println(branch.ParentItems)
//
//	if err := Models.DB.Save(&branch).Error; err != nil {
//		ReturnErr(c, err)
//	}
//
//	if err := Models.DB.Save(&branch.ParentItems).Error; err != nil {
//		ReturnErr(c, err)
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Branch Updated"})
//}

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

func FetchBranchData(c *gin.Context) {
	branch, err := getBranchByContext(c)
	if err != nil {
		ReturnErr(c, err)
	}
	if len(branch.Shifts) != 0 {
		var currentShift Models.Shift
		if err := Models.DB.Model(&Models.Shift{}).Where("branch_id = ?", branch.ID).Last(&currentShift).Error; err != nil {
			ReturnErr(c, err)
		}
		var currentEmployee Models.Employee
		if err := Models.DB.Model(&Models.Employee{}).Where("id = ?", currentShift.EmployeeID).Find(&currentEmployee).Error; err != nil {
			ReturnErr(c, err)
		}
		currentShift.Employee = currentEmployee
		branch.CurrentShift = currentShift
	}
	for shiftIndex := range branch.Shifts {
		var Employee Models.Employee
		if err := Models.DB.Model(&Models.Employee{}).Where("id = ?", branch.Shifts[shiftIndex].EmployeeID).Find(&Employee).Error; err != nil {
			ReturnErr(c, err)
		}
		branch.Shifts[shiftIndex].Employee = Employee
	}
	var ParentItems []Models.ParentItem
	for _, parentItem := range branch.ParentItems {
		var ParentItem Models.ParentItem
		var SubItems []Models.Item
		if err := Models.DB.Model(&Models.ParentItem{}).Where("id = ?", parentItem.ID).Preload("Items").Find(&ParentItem).Error; err != nil {
			ReturnErr(c, err)
		}
		for _, subItem := range ParentItem.Items {
			if !subItem.IsSold {
				SubItems = append(SubItems, subItem)
			}
		}
		ParentItem.Items = SubItems
		ParentItems = append(ParentItems, ParentItem)
	}
	var Transactions []Models.Transaction
	for _, transaction := range branch.Transactions {
		var Transaction Models.Transaction
		if err := Models.DB.Model(&Models.Transaction{}).Where("id = ?", transaction.ID).Preload("Items").Find(&Transaction).Error; err != nil {
			ReturnErr(c, err)
		}
		var ItemStructs []Models.Item
		for _, itemID := range transaction.Items {
			var Item Models.Item
			if err := Models.DB.Model(&Models.Item{}).Where("id = ?", itemID).Find(&Item).Error; err != nil {
				ReturnErr(c, err)
			}
			ItemStructs = append(ItemStructs, Item)
		}
		Transaction.ItemsStruct = ItemStructs
		Transactions = append(Transactions, Transaction)
	}
	branch.ParentItems = ParentItems
	branch.Transactions = Transactions

	var transactionsToday []Models.Transaction
	today := getCurrentFormattedDate()
	if err := Models.DB.Model(&Models.Transaction{}).Where("branch_id = ? AND DATE(created_at) = ?", branch.ID, today).Find(&transactionsToday).Error; err != nil {
		ReturnErr(c, err)
	}

	for _, transaction := range transactionsToday {
		branch.HeatMap.TotalSold += transaction.TotalCost
	}
	for _, transaction := range branch.Transactions {
		branch.TotalSold += transaction.TotalCost
	}
	branch.SoldToday = branch.HeatMap.TotalSold
	c.JSON(http.StatusOK, branch)
}
