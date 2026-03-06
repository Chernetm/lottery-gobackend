package db

import (
	"log"
	"lottery-backend/internal/config"
	"lottery-backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(mysql.Open(config.AppConfig.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	// Auto-migrate models
	err = DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Item{},
		&models.Lottery{},
		&models.Ticket{},
		&models.Coupon{},
		&models.Withdrawal{},
		&models.Payment{},
		&models.LotteryPrize{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migration completed")
}
