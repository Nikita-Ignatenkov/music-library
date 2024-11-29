package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() {
	log.Println("DEBUG: Attempting to load .env file")
	if err := godotenv.Load(); err != nil {
		log.Fatal("ERROR: Unable to load .env file:", err)
	}
	log.Println("INFO: Successfully loaded .env file")
}

func GetDatabaseURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("WARNING: DATABASE_URL is not set")
	} else {
		log.Println("DEBUG: Retrieved DATABASE_URL")
	}
	return dbURL
}
