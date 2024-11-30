package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"music-library/app/config"
	"music-library/app/models"
)

var DB *gorm.DB

func ConnectDatabase() {
	log.Println("DEBUG: Attempting to connect to the database...")
	var err error
	dsn := config.GetDatabaseURL()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("ERROR: Failed to connect to database:", err)
	}
	log.Println("INFO: Successfully connected to the database.")

	// Автомиграция
	log.Println("DEBUG: Running database migrations...")
	err = DB.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatal("ERROR: Failed to migrate database:", err)
	}
	log.Println("INFO: Database migration completed successfully.")
}
