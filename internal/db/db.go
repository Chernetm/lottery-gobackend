package db

import (
	"log"
	"lottery-backend/internal/config"
	"lottery-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.AppConfig.DBURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connection established")

	// Auto-migrate models
	// Use a two-pass approach to handle circular dependencies (e.g., Ticket <-> LotteryPrize)
	// 1. Disable FK constraints for the first pass to ensure all tables are created
	DB.Config.DisableForeignKeyConstraintWhenMigrating = true
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
		log.Fatal("Failed to migrate database (pass 1):", err)
	}

	// 2. Enable FK constraints and run again to add them properly
	DB.Config.DisableForeignKeyConstraintWhenMigrating = false
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
		log.Fatal("Failed to migrate database (pass 2):", err)
	}
	log.Println("Database migration completed")
}
