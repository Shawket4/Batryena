package Controllers

import (
	"BatrynaBackend/Models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GenerateShiftOTP(c *gin.Context) {
	const RangeLockInKilometers = 1
	var input struct {
		BranchID uint            `json:"branch_id"`
		Location Models.LatLng   `json:"lat_lng"`
		Employee Models.Employee `json:"employee"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}

	var branch Models.Branch
	if err := Models.DB.Model(&Models.Branch{}).Preload("LatLng").Where("id = ?", input.BranchID).Find(&branch).Error; err != nil {
		ReturnErr(c, err)
	}
	var Employee Models.Employee
	if err := Models.DB.Model(&Models.Employee{}).Where("name = ? AND password = ?", input.Employee.Name, input.Employee.Password).Find(&Employee).Error; err != nil {
		ReturnErr(c, err)
	}
	if Employee.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Employee Credentials"})
	}
	isInRange := CheckIfInGeographicalRange(branch.LatLng, input.Location, RangeLockInKilometers)
	if isInRange {

		var OTP Models.OTP
		OTP.BranchID = input.BranchID
		token, err := generateOTPToken(12)
		if err != nil {
			ReturnErr(c, err)
		}

		OTP.Token = token
		OTP.EmployeeID = Employee.ID
		if err := Models.DB.Model(&Models.OTP{}).Save(&OTP).Error; err != nil {
			ReturnErr(c, err)
		}

		c.JSON(http.StatusOK, gin.H{"otp": OTP.Token})
		return
	} else {
		c.JSON(http.StatusLocked, gin.H{"message": "Please Be In Range Of The Branch"})
	}
}

type SwitchShiftInput struct {
	Shift Models.Shift
	Token string `json:"token"`
}

func SwitchShift(c *gin.Context) {
	branch, err := getBranchByContext(c)
	if err != nil {
		ReturnErr(c, err)
	}
	var input SwitchShiftInput
	var currentShift Models.Shift
	var OTP Models.OTP
	var Employee Models.Employee
	if err := c.ShouldBindJSON(&input); err != nil {
		ReturnErr(c, err)
	}
	if err := Models.DB.Model(&Models.OTP{}).Where("token = ?", input.Token).Find(&OTP).Error; err != nil {
		ReturnErr(c, err)
	}
	if OTP.ID != 0 {
		if OTP.BranchID != branch.ID {
			if err := Models.DB.Model(&Models.OTP{}).Delete(&OTP).Error; err != nil {
				ReturnErr(c, err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"message": "OTP For Another Branch"})
			return
		} else {
			if err := Models.DB.Model(&Models.Employee{}).Where("id = ?", OTP.EmployeeID).Find(&Employee).Error; err != nil {
				ReturnErr(c, err)
			}
			currentTime := time.Now()
			input.Shift.BranchID = branch.ID
			input.Shift.StartedAt = currentTime
			input.Shift.BranchID = branch.ID
			input.Shift.EmployeeID = Employee.ID
			input.Shift.Employee = Employee
			if err := Models.DB.Model(&Models.Shift{}).Where("branch_id = ?", branch.ID).Last(&currentShift).Error; err != nil {
				//ReturnErr(c, err)
			}

			if currentShift.ID != 0 {
				currentShift.ClosedAt = currentTime
				currentShift.IsClosed = true
				if err := Models.DB.Save(&currentShift).Error; err != nil {
					ReturnErr(c, err)
				}
			}

			if err := Models.DB.Model(&Models.Shift{}).Create(&input.Shift).Error; err != nil {
				ReturnErr(c, err)
			}
			if err := Models.DB.Model(&Models.OTP{}).Delete(&OTP).Error; err != nil {
				ReturnErr(c, err)
			}
			c.JSON(http.StatusOK, gin.H{"message": "Shift Switched"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect OTP"})
	}
}
