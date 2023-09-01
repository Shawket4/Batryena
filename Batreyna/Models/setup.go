package Models

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {

	db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	if err != nil {
		fmt.Println("Cannot connect to database ")
		log.Fatal("connection error:", err)
	} else {
		fmt.Println("We are connected to the database ")
	}

	DB.AutoMigrate(&User{}, &Item{}, &Branch{}, &LatLng{}, &Transaction{}, &HeatMap{}, &ItemID{}, &ParentItem{})
	// DB.AutoMigrate(&DoctorWorkingHour{})
	DB.Session(&gorm.Session{FullSaveAssociations: true})
	// password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	// var user User = User{
	// 	Username:   "shawket",
	// 	Password:   string(password),
	// 	Permission: 2,
	// }
	// DB.Create(&user)
}
