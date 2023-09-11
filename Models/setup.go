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

	err = DB.AutoMigrate(&User{}, &Item{}, &Branch{}, &LatLng{}, &Transaction{}, &HeatMap{}, &ItemID{}, &ParentItem{}, &Employee{}, &Shift{}, &OTP{})
	if err != nil {
		panic(err)
	}
	DB.Session(&gorm.Session{FullSaveAssociations: true})
}
