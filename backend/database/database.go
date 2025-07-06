package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"playslate-backend/models"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	
	// Use PostgreSQL in production, SQLite for development
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		// Production: PostgreSQL
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		// Development: SQLite
		DB, err = gorm.Open(sqlite.Open("playslate.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	if err := DB.AutoMigrate(&models.User{}, &models.Subscription{}, &models.MagicLink{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func GetDB() *gorm.DB {
	return DB
}