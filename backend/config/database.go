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
	//dsn := "host=localhost user=postgres password=hainhat2003 dbname=graduation_invitation port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	dsn := "host=dpg-d4nupier433s73eeqou0-a user=gra_inv_user password=bhrUNR5HobTAZq4kDTD81GEuy3Wp9tZi dbname=gra_inv port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Show SQL queries
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")

	err = DB.AutoMigrate(
		&models.User{},
		&models.RSVP{},
		&models.Setting{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migrated successfully")

	// Seed default settings if not exist
	seedDefaultSettings()
}

func seedDefaultSettings() {
	defaultSettings := []models.Setting{
		{
			Key:         "introduction_text",
			Value:       "<p>Chào mừng bạn đến với buổi lễ tốt nghiệp!</p>",
			Description: "Nội dung giới thiệu hiển thị trên trang chủ",
		},
	}

	for _, setting := range defaultSettings {
		var existing models.Setting
		if err := DB.Where("key = ?", setting.Key).First(&existing).Error; err != nil {
			// Setting không tồn tại, tạo mới
			DB.Create(&setting)
			fmt.Printf("✅ Created default setting: %s\n", setting.Key)
		}
	}
}
