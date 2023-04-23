package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error

	// Create database connection
	db, err = gorm.Open(sqlite.Open("prince.db"), &gorm.Config{})

	if err != nil {
		log.Panicln("Failed to open database connection:", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&RepeatedMessage{})
	if err != nil {
		log.Panicln("Failed to migrate database:", err)
	}
}
