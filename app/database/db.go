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
	var err error
	dsn := config.GetDatabaseURL()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автоматически мигрируем изменения в структуре
	err = DB.AutoMigrate(&models.Song{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
