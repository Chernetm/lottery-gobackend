package main

import (
	"log"
	"lottery-backend/internal/config"
	"lottery-backend/internal/db"
	"lottery-backend/internal/handlers"
	"lottery-backend/internal/repo"
	"lottery-backend/internal/routes"
	"lottery-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Config
	config.LoadConfig()

	// Initialize DB
	db.InitDB()

	// Initialize Repositories
	userRepo := repo.NewUserRepo(db.DB)
	adminRepo := repo.NewAdminRepo(db.DB)
	itemRepo := repo.NewItemRepo(db.DB)
	lotteryRepo := repo.NewLotteryRepo(db.DB)
	ticketRepo := repo.NewTicketRepo(db.DB)
	withdrawalRepo := repo.NewWithdrawalRepo(db.DB)
	paymentRepo := repo.NewPaymentRepo(db.DB)
	couponRepo := repo.NewCouponRepo(db.DB)

	// Initialize Services
	authService := services.NewAuthService(userRepo, adminRepo)
	lotteryService := services.NewLotteryService(lotteryRepo, itemRepo)
	ticketService := services.NewTicketService(ticketRepo, userRepo, lotteryRepo, couponRepo)
	adminService := services.NewAdminService(itemRepo, lotteryRepo, userRepo, ticketRepo, couponRepo)
	withdrawalService := services.NewWithdrawalService(withdrawalRepo, userRepo)
	paymentService := services.NewPaymentService(paymentRepo, userRepo, ticketRepo, lotteryRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	lotteryHandler := handlers.NewLotteryHandler(lotteryService)
	ticketHandler := handlers.NewTicketHandler(ticketService, paymentService)
	adminHandler := handlers.NewAdminHandler(adminService)
	withdrawalHandler := handlers.NewWithdrawalHandler(withdrawalService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Initialize Gin
	r := gin.Default()

	// Set up Routes
	routes.SetupRoutes(r, authHandler, lotteryHandler, ticketHandler, adminHandler, withdrawalHandler, paymentHandler)

	// Start Server
	log.Printf("Server starting on port %s", config.AppConfig.Port)
	if err := r.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
