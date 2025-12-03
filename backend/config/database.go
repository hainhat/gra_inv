package config

import (
	"fmt"
	"graduation_invitation/backend/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=hainhat2003 dbname=graduation_invitation port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Show SQL queries
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")

	// Auto migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.RSVP{},
		&models.Setting{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migrated successfully")
}
