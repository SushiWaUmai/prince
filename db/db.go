package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	// Create data directory
	err = os.Mkdir("data", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Panicln(err)
	}

	// Create database connection
	db, err = gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{})

	if err != nil {
		log.Panicln("Failed to open database connection:", err)
	}

	err = db.AutoMigrate(&RepeatedMessage{})
	if err != nil {
		log.Panicln("Failed to migrate database:", err)
	}

	err = db.AutoMigrate(&UserPermission{})
	if err != nil {
		log.Panicln("Failed to migrate database:", err)
	}

	err = db.AutoMigrate(&MessageEvent{})
	if err != nil {
		log.Panicln("Failed to migrate database:", err)
	}

	err = db.AutoMigrate(&Alias{})
	if err != nil {
		log.Panicln("Failed to migrate database:", err)
	}
}
