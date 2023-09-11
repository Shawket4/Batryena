package Controllers

import (
	"BatrynaBackend/Models"
	"BatrynaBackend/Token"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

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
	if totalSold == 0 {
		return branches
	}
	ratio = highestSold / totalSold
	for index := range branches {
		value := branches[index].HeatMap.TotalSold / totalSold
		value /= ratio
		branches[index].HeatMap.Value = 100 * value
	}
	return branches
}

func getCurrentFormattedDate() string {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")
	return formattedDate
}

func getBranchByContext(c *gin.Context) (Models.Branch, error) {
	userId, err := Token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return Models.Branch{}, err
	}

	user, err := Models.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return Models.Branch{}, err
	}

	var branch Models.Branch
	if err := Models.DB.Model(&Models.Branch{}).Where("id = ?", user.BranchID).Preload("LatLng").Preload("ParentItems").Preload("Transactions").Preload("Shifts").Find(&branch).Error; err != nil {
		return Models.Branch{}, err
	}
	return branch, nil
}

func getUserByContext(c *gin.Context) (Models.User, error) {
	user_id, err := Token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return Models.User{}, err
	}

	user, err := Models.GetUserByID(user_id)
	return user, nil
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}
func CheckIfInGeographicalRange(Origin Models.LatLng, Point Models.LatLng, Range float64) bool {
	const earthRadiusKm = 6371.0
	// Degrees To Radians
	latOriginRad := degreesToRadians(Origin.Lat)
	longOriginRad := degreesToRadians(Origin.Lng)
	latPointRad := degreesToRadians(Point.Lat)
	longPointRad := degreesToRadians(Point.Lng)
	// Difference
	deltaLat := latPointRad - latOriginRad
	deltaLon := longPointRad - longOriginRad

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(latOriginRad)*math.Cos(latPointRad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// Calculate distance in kilometers
	distance := earthRadiusKm * c

	var output bool
	if distance <= Range {
		output = true
	} else {
		output = false
	}
	return output
}

func checkOTPTokenDuplicate(currentOTPs []Models.OTP, token string) bool {
	var output bool = false
	for _, otp := range currentOTPs {
		if otp.Token == token {
			output = true
		}
	}
	return output
}
func generateOTPToken(count int) (string, error) {
	var possibleCharacters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789#$@")
	var currentOTPs []Models.OTP
	if err := Models.DB.Model(&Models.OTP{}).Find(&currentOTPs).Error; err != nil {
		return "", err
	}

	token := make([]rune, count)
	for index := range token {
		token[index] = possibleCharacters[rand.Intn(len(possibleCharacters))]
	}
	tokenString := string(token)
	isDuplicated := checkOTPTokenDuplicate(currentOTPs, tokenString)
	if isDuplicated {
		tokenString, _ = generateOTPToken(count)
	}
	return tokenString, nil
}
