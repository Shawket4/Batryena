package Controllers

import (
	"BatrynaBackend/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterEmployee(c *gin.Context) {
	var input Models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	if err := Models.DB.Model(&Models.Employee{}).Create(&input).Error; err != nil {
		ReturnErr(c, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Employee Registered"})
}
