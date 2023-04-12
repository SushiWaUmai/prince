package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	setupDotenv()
}

func setupDotenv() {
	if _, err := os.Stat("./.env"); err == nil {
		log.Println("Found a .env file. Loading...")

		if err := godotenv.Load(); err != nil {
			log.Println("Failed to load env file")
			os.Exit(1)
		}
	}

	loadEnv()
}
