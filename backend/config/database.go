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
	dsn := "host=dpg-d4nupier433s73eeqou0-a.singapore-postgres.render.com user=gra_inv_user password=bhrUNR5HobTAZq4kDTD81GEuy3Wp9tZi dbname=gra_inv port=5432 sslmode=require"
	// dsn := "host=dpg-d4nupier433s73eeqou0-a user=gra_inv_user password=bhrUNR5HobTAZq4kDTD81GEuy3Wp9tZi dbname=gra_inv port=5432 sslmode=disable"

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
